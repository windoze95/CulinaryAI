import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from './App';
import { useParams } from 'react-router-dom';
import Footer from './Footer';

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
      {subRecipes && subRecipes.length > 0 && <h2>{mainRecipe.recipe_name}</h2>}
      <IngredientList ingredients={mainRecipe.ingredients} />
      <InstructionsList instructions={mainRecipe.instructions} />
      <p>Time to cook: {mainRecipe.time_to_cook} minutes</p>
      
      {subRecipes && subRecipes.length > 0 && (
        subRecipes.map((subRecipe, index) => (
            <div key={index}>
            <h3>{subRecipe.recipe_name}</h3>
            <IngredientList ingredients={subRecipe.ingredients} />
            <InstructionsList instructions={subRecipe.instructions} />
            <p>Time to cook: {subRecipe.time_to_cook} minutes</p>
            </div>
        ))
      )}
    </div>
  );
  
  const Recipe = () => {
    const { isAuthenticated, user } = useAuth();
    const [recipe, setRecipe] = useState(null);
    const [isGenerating, setIsGenerating] = useState(true);
  
    const { id } = useParams();
  
    const fetchRecipe = async () => {
      try {
        const response = await axios.get(`/api/v1/recipes/${id}`);
        if (response.data) {
          console.log('Recipe:', response.data.recipe);
          setRecipe(response.data.recipe);
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
    }, [isGenerating]); // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <div>
        {isGenerating ? (
          <p>Generating your recipe...<br />This may take a few minutes to complete</p>
        ) : (
          <div>
            <h1>{recipe.Title}</h1>
            {recipe && <RecipeDetail mainRecipe={recipe.FullRecipe.main_recipe} subRecipes={recipe.FullRecipe.sub_recipes} />}
            {isAuthenticated && recipe.GeneratedByUserID === user.ID && (
              <button onClick={regenerateRecipe}>Regenerate</button>
            )}
          </div>
        )}
        <Footer />
      </div>
    );
};

export default Recipe;
