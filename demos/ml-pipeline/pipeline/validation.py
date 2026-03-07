from typing import List, Dict, Tuple
import numpy as np

class PipelineValidator:
    """Validates data quality at each pipeline stage."""
    
    def __init__(self, embedding_dim: int = 128):
        self.embedding_dim = embedding_dim
        self.errors: List[str] = []
    
    def validate_embeddings(self, embeddings: np.ndarray) -> bool:
        """Check for NaN, Inf, and dimension mismatch."""
        if embeddings.shape[-1] != self.embedding_dim:
            self.errors.append(f"Dim mismatch: got {embeddings.shape[-1]}, want {self.embedding_dim}")
            return False
        if np.any(np.isnan(embeddings)):
            self.errors.append("NaN values detected in embeddings")
            return False
        if np.any(np.isinf(embeddings)):
            self.errors.append("Inf values detected in embeddings")
            return False
        return True
    
    def validate_feature_distribution(self, 
                                       features: np.ndarray,
                                       expected_mean: float = 0.0,
                                       tolerance: float = 0.5) -> Tuple[bool, str]:
        """Detect train/serve skew via distribution check."""
        actual_mean = float(np.mean(features))
        drift = abs(actual_mean - expected_mean)
        if drift > tolerance:
            msg = f"Distribution drift detected: mean={actual_mean:.3f}, expected={expected_mean:.3f}"
            return False, msg
        return True, "ok"
