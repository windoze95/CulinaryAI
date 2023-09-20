import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Paper from '@material-ui/core/Paper';
import Grid from '@material-ui/core/Grid';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import swal from 'sweetalert';
import axios from 'axios';

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
    avatar: {
      margin: theme.spacing(0.5),   // Reduced margin
      backgroundColor: theme.palette.secondary.main,
    },
    form: {
      width: '80%',  // Reduced width
      marginTop: theme.spacing(1),
    },
    submit: {
      margin: theme.spacing(1, 0, 1),  // Reduced margin
    },
  }));

  async function loginUser(credentials) {
    try {
      const response = await axios.post('/api/v1/users/login', credentials, {
        headers: {
          'Content-Type': 'application/json'
        }
      });
      return response.data;
    } catch (error) {
      return { message: 'Login failed' };
    }
  }

export default function Signin() {
  const classes = useStyles();
  const [username, setUserName] = useState('');
  const [password, setPassword] = useState('');

  // const navigate = useNavigate();

  const handleSubmit = async e => {
    e.preventDefault();
    const response = await loginUser({
      username,
      password
    });
    if ('accessToken' in response) {
      swal("Success", response.message, "success", {
        buttons: false,
        timer: 2000,
      })
      .then((value) => {
        localStorage.setItem('accessToken', response['accessToken']);
        localStorage.setItem('user', JSON.stringify(response['user']));
        window.location.href = "/";
        // navigate("/profile");
      });
    } else {
      swal("Failed", response.message, "error");
    }
  }

  return (
    <Grid container className={classes.root} justifyContent="center"> {/* Added justifyContent */}
    <CssBaseline />
    <Grid item xs={12} md={7} component={Paper} elevation={6} square>
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Sign in
          </Typography>
          <form className={classes.form} noValidate onSubmit={handleSubmit}>
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="username"
              name="username"
              label="Username"
              onChange={e => setUserName(e.target.value)}
            />
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="password"
              name="password"
              label="Password"
              type="password"
              onChange={e => setPassword(e.target.value)}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
            >
              Sign In
            </Button>
            <Typography variant="body2" align="center">Don't have an account?</Typography>
            <Link to="/register" style={{ textDecoration: 'none' }}>
              <Typography variant="body2" align="center">
                Register
              </Typography>
            </Link>
          </form>
        </div>
      </Grid>
    </Grid>
  );
}