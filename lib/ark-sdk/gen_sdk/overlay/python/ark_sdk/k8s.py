"""Kubernetes utilities and client initialization."""
import functools
import logging
import os
from functools import lru_cache

from kubernetes import config
from kubernetes.config.config_exception import ConfigException
from kubernetes_asyncio import config as async_config

logger = logging.getLogger(__name__)

NS_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

@functools.lru_cache(maxsize=1)
def get_namespace():
    ns = 'default'
    if is_k8s() and os.path.isfile(NS_PATH):
        with open(NS_PATH) as f:
            ns = f.read().strip()
    logger.info(f"Using namespace {ns}")
    return ns

@functools.lru_cache(maxsize=1)
def is_k8s():
    return os.path.isfile(NS_PATH)

@lru_cache(maxsize=1)
def _init_k8s():
    """Initialize Kubernetes client configuration."""
    try:
        # Load kubeconfig from default location (~/.kube/config)
        config.load_kube_config()
        logger.info("Loaded kubeconfig from default location (probably dev mode)")
        
        # Log the current context for debugging
        contexts, active_context = config.list_kube_config_contexts()
        if active_context:
            logger.info(f"Active context: {active_context['name']}")
            
    except ConfigException:
        try:
            # Try to load in-cluster config if running inside a pod
            config.load_incluster_config()
            logger.info("Loaded in-cluster config")
        except ConfigException as e:
            logger.error(f"Failed to load any Kubernetes config: {e}")
            raise


async def init_k8s():
    """Initialize Kubernetes async client configuration by wrapping sync init."""
    # First ensure sync config is loaded in case we need it
    _init_k8s()
    
    # Then load the async config using the same method
    try:
        await async_config.load_kube_config()
    except:
        # If that fails, try in-cluster config
        async_config.load_incluster_config()