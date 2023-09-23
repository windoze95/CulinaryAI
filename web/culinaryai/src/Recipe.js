import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from './App';

const IngredientList = ({ ingredients }) => (
    <ul>
      {ingredients.map((ingredient, index) => (
        <li key={index}>
          {ingredient.amount} {ingredient.unit} of {ingredient.name}
        </li>
      ))}
    </ul>
);

const InstructionsList = ({ instructions }) => (
    <ol>
        {instructions.map((instruction, index) => (
        <li key={index}>{instruction}</li>
        ))}
    </ol>
);

const RecipeDetail = ({ mainRecipe, subRecipes }) => (
    <div>
        <h2>Main Recipe</h2>
        <IngredientList ingredients={mainRecipe.ingredients} />
        <InstructionsList instructions={mainRecipe.instructions} />
        <p>Time to cook: {mainRecipe.timeToCook} minutes</p>
        
        {subRecipes.map((subRecipe, index) => (
        <div key={index}>
            <h3>Sub Recipe {index + 1}</h3>
            <IngredientList ingredients={subRecipe.ingredients} />
            <InstructionsList instructions={subRecipe.instructions} />
            <p>Time to cook: {subRecipe.timeToCook} minutes</p>
        </div>
        ))}
    </div>
);

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
                <h1>{recipe.Title}</h1>
                {/* Display your recipe here */}
                {recipe && <RecipeDetail mainRecipe={recipe.FullRecipe.MainRecipe} subRecipes={recipe.FullRecipe.SubRecipes} />}
                {isAuthenticated && recipe.GeneratedByUserID === user.ID && (
                    <button onClick={regenerateRecipe}>Regenerate</button>
                )}
            </div>
        )}
      </div>
    );
};

export default Recipe;
