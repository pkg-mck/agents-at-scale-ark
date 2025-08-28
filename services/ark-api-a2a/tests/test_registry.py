import unittest
from unittest import IsolatedAsyncioTestCase
from unittest.mock import AsyncMock, MagicMock, patch

from a2agw.registry import AgentRegistry, ark_to_agent_card


class TestAgentRegistry(IsolatedAsyncioTestCase):
    def setUp(self):
        """Set up test fixtures"""
        # Create mock ARK agent
        self.mock_ark_agent = MagicMock()
        self.mock_ark_agent.metadata = {"name": "test-agent"}
        self.mock_ark_agent.spec.description = "Test agent description"

        # Create mock ARK client
        self.mock_client = AsyncMock()
        self.mock_context = AsyncMock()
        self.mock_context.__aenter__.return_value = self.mock_client
        self.mock_context.__aexit__.return_value = None

        # Patch the with_ark_client function
        self.patcher = patch("a2agw.registry.with_ark_client")
        self.mock_with_ark_client = self.patcher.start()
        self.mock_with_ark_client.return_value = self.mock_context

    def tearDown(self):
        """Clean up after tests"""
        self.patcher.stop()

    async def test_registry_initialization(self):
        registry = AgentRegistry(namespace="test-namespace")
        self.assertEqual(registry._namespace, "test-namespace")

    async def test_list_agents(self):
        # Setup mock agents
        agent1 = MagicMock()
        agent1.metadata = {"name": "agent-1"}
        agent1.spec.description = "Agent 1 description"

        agent2 = MagicMock()
        agent2.metadata = {"name": "agent-2"}
        agent2.spec.description = "Agent 2 description"

        self.mock_client.agents.a_list.return_value = [agent1, agent2]

        registry = AgentRegistry(namespace="test-namespace")
        agent_cards = await registry.list_agents()

        self.assertEqual(len(agent_cards), 2)
        self.assertEqual(agent_cards[0].name, "agent-1")
        self.assertEqual(agent_cards[0].description, "Agent 1 description")
        self.assertEqual(agent_cards[1].name, "agent-2")
        self.assertEqual(agent_cards[1].description, "Agent 2 description")

    async def test_get_agent(self):
        self.mock_client.agents.a_get.return_value = self.mock_ark_agent

        registry = AgentRegistry(namespace="test-namespace")
        agent_card = await registry.get_agent("test-agent")

        self.assertIsNotNone(agent_card)
        self.assertEqual(agent_card.name, "test-agent")
        self.assertEqual(agent_card.description, "Test agent description")
        self.mock_client.agents.a_get.assert_called_once_with("test-agent")

    async def test_get_nonexistent_agent(self):
        self.mock_client.agents.a_get.side_effect = Exception("Agent not found")

        registry = AgentRegistry(namespace="test-namespace")

        with self.assertRaisesRegex(Exception, "Agent not found"):
            await registry.get_agent("nonexistent")

    async def test_find_agents_by_capability(self):
        # Create agents with different capabilities
        agent1 = MagicMock()
        agent1.metadata = {"name": "data-agent"}
        agent1.spec.description = "Data processing agent"

        agent2 = MagicMock()
        agent2.metadata = {"name": "task-agent"}
        agent2.spec.description = "Task management agent"

        self.mock_client.agents.a_list.return_value = [agent1, agent2]

        registry = AgentRegistry(namespace="test-namespace")

        # Currently, ark_to_agent_card creates agents with default "General" skill
        # So this test should return agents with matching skill names
        agents = await registry.find_agents_by_capability("General")
        self.assertEqual(len(agents), 2)  # Both agents have default "General" skill

        agents = await registry.find_agents_by_capability("nonexistent")
        self.assertEqual(len(agents), 0)

    async def test_ark_to_agent_card_conversion(self):
        # Test the ark_to_agent_card conversion function
        ark_agent = MagicMock()
        ark_agent.metadata = {"name": "converted-agent"}
        ark_agent.spec.description = "Converted agent description"

        agent_card = ark_to_agent_card(ark_agent)

        self.assertEqual(agent_card.name, "converted-agent")
        self.assertEqual(agent_card.description, "Converted agent description")
        self.assertEqual(len(agent_card.skills), 1)  # Default skill
        self.assertEqual(agent_card.skills[0].name, "General")
        self.assertEqual(agent_card.url, "http://localhost:7184/agent/converted-agent/")
        self.assertEqual(agent_card.version, "1.0.0")

    async def test_ark_to_agent_card_no_description(self):
        # Test conversion when description is None
        ark_agent = MagicMock()
        ark_agent.metadata = {"name": "no-desc-agent"}
        ark_agent.spec.description = None

        agent_card = ark_to_agent_card(ark_agent)

        self.assertEqual(agent_card.name, "no-desc-agent")
        self.assertEqual(agent_card.description, "No description")



if __name__ == "__main__":
    unittest.main()
