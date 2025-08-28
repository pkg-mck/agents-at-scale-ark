#!/usr/bin/env node
import fs from 'fs';
import path, { dirname } from 'path';
import { fileURLToPath } from 'url';

import axios from 'axios';
import { execa } from 'execa';
import { Text } from 'ink';
import open from 'open';
import { useState, useEffect } from 'react';
import YAML from 'yaml';

import {
  DEFAULT_ADDRESS_ARK_API,
  DEFAULT_ARK_DASHBOARD_URL,
} from '../lib/consts.js';

const DashboardCLI = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [status, setStatus] = useState<'Running' | 'Done' | 'Failed'>(
    'Running'
  );
  const [dotsIndex, setDotsIndex] = useState(0);
  const dotFrames = ['.', '..', '...'];

  const __filename = fileURLToPath(import.meta.url);
  const __dirname = dirname(__filename);
  const rootDir = path.resolve(__dirname, '../../../');

  const steps = [
    {
      label: 'Building backend',
      path: `${rootDir}/ark-api`,
      checkUrl: `${DEFAULT_ADDRESS_ARK_API}/health`,
    },
    {
      label: 'Starting frontend',
      path: `${rootDir}/ark-dashboard`,
      checkUrl: DEFAULT_ARK_DASHBOARD_URL,
    },
  ];

  async function isServiceUp(url: string): Promise<boolean> {
    try {
      const res = await axios.get(url);
      return res.status >= 200 && res.status < 400;
    } catch {
      return false;
    }
  }

  async function runCommands(commands: string[], cwd: string) {
    for (const cmd of commands) {
      // Run each command in the specified working directory
      await execa(cmd, {
        cwd: cwd,
        stdio: ['ignore', 'pipe', 'pipe'],
        shell: true,
      });
    }
  }

  function loadServiceManifest(servicePath: string) {
    const manifestPath = path.join(servicePath, 'manifest.yaml');
    const raw = fs.readFileSync(manifestPath, 'utf-8');
    return YAML.parse(raw);
  }

  async function runStep(stepIndex: number) {
    const step = steps[stepIndex];
    const manifest = loadServiceManifest(step.path);
    // 1. Check if service is already up
    const alreadyUp = await isServiceUp(step.checkUrl);
    if (alreadyUp) {
      setCurrentStep((prev) => prev + 1);
      return;
    }

    // 2. Install the dependencies
    if (manifest.commands?.install) {
      try {
        await runCommands(manifest.commands.install, step.path);
      } catch (err) {
        setStatus('Failed');
        console.error(`Error running install commands for ${step.label}`, err);
        process.exit(1);
      }
    }

    // 3. Start the service
    const startCommand = manifest.commands?.dev;
    if (startCommand) {
      const devCommandStr = Array.isArray(startCommand)
        ? startCommand.join(' && ')
        : startCommand;
      const serviceProcess = execa(devCommandStr, {
        cwd: step.path,
        stdio: ['ignore', 'pipe', 'pipe'],
        shell: true,
      });

      // 4. Poll until service is ready
      return new Promise<void>((resolve, reject) => {
        let attempts = 0;
        const maxAttempts = 30;
        const interval = 3000;

        const poll = async () => {
          attempts++;
          if (await isServiceUp(step.checkUrl)) {
            setCurrentStep((prev) => prev + 1);
            resolve();
            return;
          }

          if (attempts >= maxAttempts) {
            reject(new Error(`${step.label} did not start in time.`));
            return;
          }

          setTimeout(poll, interval);
        };

        poll();

        serviceProcess.catch((err: any) => reject(err));
      });
    }
  }

  async function startDashboardServices() {
    try {
      for (let i = 0; i < steps.length; i++) {
        await runStep(i);
      }
      setStatus('Done');
      await open(steps[steps.length - 1].checkUrl);
    } catch (err) {
      setStatus('Failed');
      console.error('Error:', err);
      process.exit(1);
    }
  }

  useEffect(() => {
    startDashboardServices();
  }, []);

  useEffect(() => {
    if (status !== 'Running') {
      return;
    }
    const interval = setInterval(() => {
      setDotsIndex((prev) => (prev + 1) % dotFrames.length);
    }, 300);
    return () => clearInterval(interval);
  }, [status]);

  return (
    <>
      <Text color="yellow">⏳ Starting Dashboard</Text>
      {steps.map((step, i) => {
        if (i < currentStep) {
          return (
            <Text key={i} color="green">
              ✅ Step {i + 1}: {step.label}
            </Text>
          );
        } else if (i === currentStep && status === 'Running') {
          return (
            <Text key={i}>
              Step {i + 1}: {step.label}
              {dotFrames[dotsIndex]}
            </Text>
          );
        } else {
          return null;
        }
      })}

      {status === 'Done' && (
        <Text color="green">🎉 Dashboard is up and running!</Text>
      )}
      {status === 'Failed' && (
        <Text color="red">❌ Failed to start Dashboard</Text>
      )}
    </>
  );
};

export default DashboardCLI;
