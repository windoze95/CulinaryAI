import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import ReCAPTCHA from "react-google-recaptcha";
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

const validatePassword = (password) => {
  if (password.length < 8) return 'Password must be at least 8 characters long';
  if (!/[A-Z]/.test(password)) return 'Password must contain at least one uppercase letter';
  if (!/[a-z]/.test(password)) return 'Password must contain at least one lowercase letter';
  if (!/\d/.test(password)) return 'Password must contain at least one digit';
  if (!/[!@#$%^&*]/.test(password)) return 'Password must contain at least one special character';
  return null;
};

export default function Register() {
  const classes = useStyles();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordError, setPasswordError] = useState(null);
  const [isVerified, setIsVerified] = useState(false);

  const handlePasswordChange = (e) => {
    const newPass = e.target.value;
    setPassword(newPass);
    setPasswordError(validatePassword(newPass));
  };

  const handleCaptchaResponse = (value) => {
    if (value) {
      setIsVerified(true);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
  
    // Check if passwords match
    if (password !== confirmPassword) {
      swal("Failed", "Passwords do not match", "error");
      return;
    }
  
    // Check if password meets requirements
    if (passwordError) {
      swal("Failed", passwordError, "error");
      return;
    }

    if (!isVerified) {
      swal("Failed", "Please verify that you are a human", "error");
      return;
    }
  
    try {
      const response = await axios.post('/api/v1/users/register', {
        username,
        email,
        password,
        isVerified
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      });
  
      // Assuming the API responds with a message when registration is successful
      if (response.data && response.data.message) {
        swal("Success", response.data.message, "success", {
          buttons: false,
          timer: 2000,
        })
        .then((value) => {
          // Redirect to the login page
          window.location.href = "/signin";
        });
      } else {
        // Handle case where API response is not as expected
        swal("Failed", "Registration failed", "error");
      }
  
    } catch (error) {
      // Handle API call errors
      swal("Failed", "Registration failed", "error");
    }
  };
  

  return (
    <Grid container className={classes.root} justifyContent="center">
      <CssBaseline />
      <Grid item xs={12} md={7} component={Paper} elevation={6} square>
        <div className={classes.paper}>
          <Avatar className={classes.avatar}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Register
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
              onChange={e => setUsername(e.target.value)}
            />
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="email"
              name="email"
              label="Email Address"
              onChange={e => setEmail(e.target.value)}
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
              onChange={handlePasswordChange}
              error={Boolean(passwordError)}
              helperText={passwordError}
            />
            <TextField
              variant="outlined"
              margin="normal"
              required
              fullWidth
              id="confirmPassword"
              name="confirmPassword"
              label="Confirm Password"
              type="password"
              onChange={e => setConfirmPassword(e.target.value)}
              error={Boolean(password !== confirmPassword)}
              helperText={password !== confirmPassword ? "Passwords do not match" : ""}
            />
            <ReCAPTCHA
              sitekey="6LeXecgnAAAAAHpoDsyzOZ4Zl9J7saMqaosXZh2T"
              onChange={handleCaptchaResponse}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
              disabled={!isVerified} // Disable the button until verified
            >
              Register
            </Button>
            <Typography variant="body2" align="center">Already have an account?</Typography>
            <Link to="/signin" style={{ textDecoration: 'none' }}>
              <Typography variant="body2" align="center">
                Sign In
              </Typography>
            </Link>
          </form>
        </div>
      </Grid>
    </Grid>
  );
}
