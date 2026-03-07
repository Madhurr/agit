import numpy as np
from typing import List, Dict

def compute_user_embeddings(user_history: List[Dict]) -> np.ndarray:
    """Compute dense user embeddings from interaction history."""
    if not user_history:
        return np.zeros(128)
    
    # Weighted average of item embeddings by recency
    weights = np.exp(-0.1 * np.arange(len(user_history)))
    weights /= weights.sum()
    
    embeddings = np.array([item.get('embedding', np.zeros(128)) 
                          for item in user_history])
    return np.average(embeddings, axis=0, weights=weights)

def normalize_features(features: np.ndarray) -> np.ndarray:
    """L2 normalize feature vectors."""
    norms = np.linalg.norm(features, axis=-1, keepdims=True)
    return features / (norms + 1e-8)

def compute_item_popularity(item_ids: List[str], 
                            interaction_counts: Dict[str, int]) -> np.ndarray:
    """Log-normalized popularity scores to reduce popularity bias."""
    counts = np.array([interaction_counts.get(iid, 0) for iid in item_ids])
    return np.log1p(counts) / np.log1p(counts.max() + 1)
