import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from './App';
import { useLocation, useParams } from 'react-router-dom';
import Footer from './Footer';
import { makeStyles } from '@material-ui/core/styles';
import LogoSvg from './logo.svg';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { Helmet } from 'react-helmet';

const useStyles = makeStyles((theme) => ({
    root: {
      height: '100vh',
    },
    recipeImage: {
        width: '100%',  // or you can set a specific width
        aspectRatio: 1, // maintain aspect ratio
        resizeMode: 'contain', // or 'cover'
    },
    paper: {
      margin: theme.spacing(4, 2),  // Reduced margin
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
    },
  }));

  const IngredientList = ({ ingredients }) => (
    <ul>
      {ingredients.map((ingredient, index) => (
        <li key={index}>
          {(ingredient.unit.toLowerCase() === "pieces" || ingredient.unit.toLowerCase() === "piece")
            ? `${ingredient.amount} ${ingredient.name}`
            : `${ingredient.amount} ${ingredient.unit} of ${ingredient.name}`}
        </li>
      ))}
    </ul>
  );  

// const IngredientList = ({ ingredients }) => (
//     <ul>
//       {ingredients.map((ingredient, index) => (
//         <li key={index}>
//           {ingredient.amount} {ingredient.unit} of {ingredient.name}
//         </li>
//       ))}
//     </ul>
//   );
  
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
    const classes = useStyles();
    const { isAuthenticated, user } = useAuth();
    const [recipe, setRecipe] = useState(null);
    const location = useLocation();
    const [isGenerating, setIsGenerating] = useState(false);

    useEffect(() => {
        // Check if the user came from the "/generate" route
        if (location.state && location.state.from === '/generate') {
          setIsGenerating(true);
        }
    }, [location]); // eslint-disable-line react-hooks/exhaustive-deps
  
    const { id } = useParams();
  
    const fetchRecipe = async () => {
      try {
        const response = await axios.get(`/api/v1/recipes/${id}`);
        if (response.data && response.data.recipe) {
            setRecipe(response.data.recipe);
        
            // Set isGenerating based on GenerationComplete property
            if (response.data.recipe.GenerationComplete === undefined || 
                response.data.recipe.GenerationComplete === false) {
                setIsGenerating(true);
            } else {
                setIsGenerating(false);
            }
        }
        // if (response.data) {
        //   setRecipe(response.data.recipe);
        //   setIsGenerating(!response.data.recipe.GenerationComplete);
        // }
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
      <Grid container className={classes.root} justifyContent="center">
        <Grid item xs={12} md={7} component={Paper} elevation={6} square>
          <div className={classes.paper}>
            <Helmet>
                <title>{recipe ? recipe.Title : 'Loading...'}</title>
                <meta property="og:title" content={recipe ? recipe.Title : 'Loading...'} />
                <meta property="og:image" content={recipe ? recipe.ImageURL : 'default-image-url'} />
                {/* <meta property="og:description" content="Your description here" /> */}
                <meta property="og:url" content={window.location.href} />
            </Helmet>
            {isGenerating ? (
            <p>Generating your recipe...<br />
            <div className="loading-div">
              <img src={LogoSvg} className="breathe-logo" alt="Logo" />
            </div><br />
            This may take a few minutes to complete</p>
            ) : (
            <div>
                <h1>{recipe ? recipe.Title : 'Loading...'}</h1>
                <img src={recipe ? recipe.ImageURL : 'default-image-url'} alt={recipe ? recipe.Title : 'Loading...'} className={classes.recipeImage} />
                {recipe && <RecipeDetail mainRecipe={recipe.FullRecipe.main_recipe} subRecipes={recipe.FullRecipe.sub_recipes} />}
                {isAuthenticated && recipe && recipe.GeneratedByUserID === user.ID && (
                <button onClick={regenerateRecipe}>Regenerate</button>
                )}
            </div>
            )}
          </div>
          <Footer />
        </Grid>
      </Grid>
    );
};

export default Recipe;
