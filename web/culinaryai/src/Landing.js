import React from 'react';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import logo from './logo.svg';
import GithubLinkButton from './GithubLinkButton';
import Footer from './Footer';

const useStyles = makeStyles((theme) => ({
    root: {
      height: '100vh',
    },
    image: {
      // Removed background image
    },
    paper: {
      margin: theme.spacing(4, 2),  // Reduced margin
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
    },
    logo: {
      width: '200px',
      height: '200px',
      marginBottom: theme.spacing(2),
    },
  }));

export default function Landing() {
  const classes = useStyles();

  return (
    <Grid container className={classes.root} justifyContent="center">
      <Grid item xs={12} md={7} component={Paper} elevation={6} square>
        <div className={classes.paper}>
          <img src={logo} className={classes.logo} alt="CulinaryAI™" />™
          <Typography component="h1" variant="h5">
            iOS and Android app coming soon!
          </Typography>
          <GithubLinkButton />
          <Typography variant="body1" align="center">
            <br />
            ** The name: <strong>CulinaryAI™</strong> in the Apple App Store is currently subject to a naming despute.<br />
            <br />
            Note to Apple: This domain is owned and operated by Juliano DiCesare, registered Apple Developer. I assert common law rights to the use of the wordmark 'CulinaryAI™' in commerce. The trademark is also submitted and currently under review for official registration.<br />
          </Typography>
        </div>
        <Footer />
      </Grid>
    </Grid>
  );
}
