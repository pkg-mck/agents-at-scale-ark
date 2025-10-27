import { useAtom } from 'jotai';
import { RESET } from 'jotai/utils';
import React, { Fragment, useCallback, useEffect, useState } from 'react';

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';

import { experimentalFeatureGroups } from './experimental-features';
import type { ExperimentalFeature } from './types';

const EXPERIMENTAL_MODAL_KEYBOARD_SHORTCUT = 'e';

type ExperimentalFeatureToggleProps = {
  feature: ExperimentalFeature;
};

function ExperimentalFeatureToggle({
  feature,
}: ExperimentalFeatureToggleProps) {
  const [atomValue, setAtom] = useAtom(feature.atom);

  const toggleAtomValue = useCallback(() => {
    setAtom(prev => (prev ? RESET : true));
  }, [setAtom]);

  return (
    <div className="flex flex-row items-center justify-between">
      <div className="space-y-0.5">
        <Label>{feature.feature}</Label>
        {feature.description && (
          <div className="text-muted-foreground text-sm">
            {feature.description}
          </div>
        )}
      </div>
      <Switch checked={atomValue} onCheckedChange={toggleAtomValue} />
    </div>
  );
}

export function ExperimentalFeaturesDialog() {
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const toggleModal = useCallback(() => {
    setIsDialogOpen(prev => !prev);
  }, []);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (
        event.key === EXPERIMENTAL_MODAL_KEYBOARD_SHORTCUT &&
        (event.metaKey || event.ctrlKey)
      ) {
        event.preventDefault();
        toggleModal();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [toggleModal]);

  return (
    <Dialog open={isDialogOpen} onOpenChange={toggleModal}>
      <DialogContent
        className="max-h-[90vh] overflow-y-auto sm:max-w-2xl"
        onOpenAutoFocus={e => e.preventDefault()}>
        <DialogHeader>
          <DialogTitle>Experimental features</DialogTitle>
          <DialogDescription>Enable experimental features</DialogDescription>
        </DialogHeader>
        <div className="space-y-6 px-2 py-4">
          {experimentalFeatureGroups.map(
            ({ groupKey, groupLabel, features }) => (
              <section key={groupKey} className="space-y-2">
                {groupLabel && (
                  <Label className="text-base font-bold">{groupLabel}</Label>
                )}
                <div>
                  {features.map((feature, index) => (
                    <Fragment key={feature.feature}>
                      {index !== 0 && <Separator />}
                      <ExperimentalFeatureToggle feature={feature} />
                    </Fragment>
                  ))}
                </div>
              </section>
            ),
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
