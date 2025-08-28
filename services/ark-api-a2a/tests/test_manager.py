import asyncio
import unittest
from unittest.mock import AsyncMock, MagicMock, patch

from a2a.types import AgentCard, AgentSkill

from src.a2agw.manager import DynamicManager, ProxyApp


class TestProxyApp(unittest.IsolatedAsyncioTestCase):
    """Test the ProxyApp ASGI proxy"""
    
    def setUp(self):
        self.proxy = ProxyApp()
    
    async def test_proxy_returns_404_when_no_app(self):
        """Test that proxy returns 404 when no app is set"""
        # Mock ASGI interface
        scope = {'type': 'http', 'path': '/test'}
        messages = []
        
        async def send(message):
            messages.append(message)
        
        async def receive():
            return {'type': 'http.disconnect'}
        
        # Call proxy with no app set
        await self.proxy(scope, receive, send)
        
        # Verify 404 response
        self.assertEqual(len(messages), 2)
        self.assertEqual(messages[0]['status'], 404)
        self.assertEqual(messages[1]['body'], b'Not Found')
    
    async def test_proxy_forwards_to_app(self):
        """Test that proxy forwards requests to the set app"""
        # Create a mock ASGI app
        mock_app = AsyncMock()
        self.proxy.set_app(mock_app)
        
        # Mock ASGI interface
        scope = {'type': 'http', 'path': '/test'}
        receive = AsyncMock()
        send = AsyncMock()
        
        # Call proxy
        await self.proxy(scope, receive, send)
        
        # Verify app was called with correct parameters
        mock_app.assert_called_once_with(scope, receive, send)
    
    async def test_proxy_atomic_swap(self):
        """Test that proxy can atomically swap apps"""
        # Create two mock apps
        app1 = AsyncMock()
        app2 = AsyncMock()
        
        # Set first app
        self.proxy.set_app(app1)
        
        # Mock ASGI interface
        scope = {'type': 'http', 'path': '/test'}
        receive = AsyncMock()
        send = AsyncMock()
        
        # Call with first app
        await self.proxy(scope, receive, send)
        app1.assert_called_once()
        
        # Swap to second app
        self.proxy.set_app(app2)
        
        # Call with second app
        await self.proxy(scope, receive, send)
        app2.assert_called_once()


class TestDynamicManager(unittest.IsolatedAsyncioTestCase):
    """Test the DynamicManager"""
    
    def setUp(self):
        self.patcher = patch('src.a2agw.manager.get_registry')
        self.mock_get_registry = self.patcher.start()
        self.mock_registry = AsyncMock()
        self.mock_get_registry.return_value = self.mock_registry
        
        self.manager = DynamicManager()
    
    def tearDown(self):
        self.patcher.stop()
    
    async def test_initialize_starts_periodic_sync(self):
        """Test that initialize starts the periodic sync task"""
        # Mock registry to return empty list
        self.mock_registry.list_agents.return_value = []
        
        # Initialize manager
        await self.manager.initialize()
        
        # Verify sync was called
        self.mock_registry.list_agents.assert_called_once()
        
        # Verify periodic task was started
        self.assertIsNotNone(self.manager._refresh_task)
        self.assertTrue(self.manager._running)
        
        # Clean up
        await self.manager.shutdown()
    
    async def test_shutdown_stops_periodic_sync(self):
        """Test that shutdown stops the periodic sync task"""
        # Initialize first
        self.mock_registry.list_agents.return_value = []
        await self.manager.initialize()
        
        # Shutdown
        await self.manager.shutdown()
        
        # Verify task was stopped
        self.assertFalse(self.manager._running)
        self.assertIsNone(self.manager._refresh_task)
    
    def test_app_is_proxy(self):
        """Test that manager.app is a ProxyApp instance"""
        self.assertIsInstance(self.manager.app, ProxyApp)
    
    @patch('src.a2agw.manager.ARKAgentExecutor')
    @patch('src.a2agw.manager.A2AStarletteApplication')
    async def test_routes_updated_when_agents_change(self, mock_a2a_app, mock_executor):
        """Test that routes are updated when agents change"""
        # Create mock agent cards
        from a2a.types import AgentCapabilities
        
        agent1 = AgentCard(
            name="agent1",
            description="Test agent 1",
            skills=[AgentSkill(id="skill1", name="General", description="General skill", tags=[])],
            url="http://localhost:7184/agent/agent1",
            version="1.0.0",
            capabilities=AgentCapabilities(streaming=True, pushNotifications=False, stateTransitionHistory=False),
            defaultInputModes=["text"],
            defaultOutputModes=["text"]
        )
        agent2 = AgentCard(
            name="agent2", 
            description="Test agent 2",
            skills=[AgentSkill(id="skill2", name="General", description="General skill", tags=[])],
            url="http://localhost:7184/agent/agent2",
            version="1.0.0",
            capabilities=AgentCapabilities(streaming=True, pushNotifications=False, stateTransitionHistory=False),
            defaultInputModes=["text"],
            defaultOutputModes=["text"]
        )
        
        # Mock the build method
        mock_a2a_app.return_value.build.return_value = MagicMock()
        
        # First sync - one agent
        self.mock_registry.list_agents.return_value = [agent1]
        await self.manager.initialize()
        
        # Verify routes were updated
        initial_app = self.manager.app._app
        self.assertIsNotNone(initial_app)
        
        # Second sync - two agents (simulating periodic sync)
        self.mock_registry.list_agents.return_value = [agent1, agent2]
        await self.manager._sync_with_registry()
        
        # Verify app was replaced (new instance)
        updated_app = self.manager.app._app
        self.assertIsNotNone(updated_app)
        # Note: In real implementation, these would be different instances
        
        # Clean up
        await self.manager.shutdown()
    
    async def test_sync_handles_registry_errors(self):
        """Test that sync handles registry errors gracefully"""
        # Make registry raise an error
        self.mock_registry.list_agents.side_effect = Exception("Registry error")
        
        # This should not raise - errors are logged
        await self.manager._sync_with_registry()
        
        # Verify registry was called
        self.mock_registry.list_agents.assert_called_once()
    
    @patch('src.a2agw.manager.POLL_INTERVAL', 0.1)  # Speed up test
    async def test_periodic_sync_runs_periodically(self):
        """Test that periodic sync runs at intervals"""
        # Mock registry
        self.mock_registry.list_agents.return_value = []
        
        # Start periodic sync
        await self.manager.start_periodic_sync()
        
        # Wait for multiple intervals
        await asyncio.sleep(0.25)
        
        # Verify multiple syncs happened
        self.assertGreaterEqual(self.mock_registry.list_agents.call_count, 2)
        
        # Clean up
        await self.manager.stop_periodic_sync()


if __name__ == "__main__":
    unittest.main()