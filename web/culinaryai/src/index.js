import React from 'react';
import ReactDOM from 'react-dom/client';
import 'materialize-css';
import 'materialize-css/dist/css/materialize.min.css';
import 'material-icons/iconfont/material-icons.css';
import './index.css';
import App from './App';
import { LoadingProvider } from './LoadingContext';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <LoadingProvider> {/* Wrap App with LoadingProvider */}
      <App />
    </LoadingProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
