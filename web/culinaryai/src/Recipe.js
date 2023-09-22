import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from './App';

const Recipe = ({ match }) => {
    const { isAuthenticated, user } = useAuth();
    const [recipe, setRecipe] = useState(null);
    const [isGenerating, setIsGenerating] = useState(true);

    const fetchRecipe = async () => {
        try {
        const response = await axios.get(`/api/v1/recipes/${match.params.id}`);
        if (response.data) {
            setRecipe(response.data);
            setIsGenerating(!response.data.recipe.GenerationComplete);
        }
        } catch (error) {
            console.error('Error fetching recipe:', error);
        }
    };

    const regenerateRecipe = async () => {
        // Logic to regenerate the recipe
    };

    useEffect(() => {
        fetchRecipe();

        const interval = setInterval(() => {
        if (isGenerating) {
            fetchRecipe();
        }
        }, 5000); // Poll every 5 seconds

        return () => clearInterval(interval);
    }, [isGenerating]);

    return (
        <div>
        {isGenerating ? (
            <p>Generating your recipe...</p>
        ) : (
            <div>
                {/* Display your recipe here */}
                {isAuthenticated && recipe.GeneratedByUserID === user.ID && (
                <button onClick={regenerateRecipe}>Regenerate</button>
                )}
            </div>
        )}
      </div>
    );
};

export default Recipe;
