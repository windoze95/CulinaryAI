import React from 'react';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
    footer: {
      marginTop: theme.spacing(4),
      padding: theme.spacing(2),
      backgroundColor: theme.palette.grey[200],
    },
    text: {
      textAlign: 'center',
    }
  }));

export default function Footer() {
  const classes = useStyles();

  return (
    <Grid container className={classes.footer}>
      <Grid item xs={12}>
        <Typography variant="body2" className={classes.text}>
          &copy; 2023 CulinaryAI™. All rights reserved.
        </Typography>
      </Grid>
    </Grid>
  );
}
