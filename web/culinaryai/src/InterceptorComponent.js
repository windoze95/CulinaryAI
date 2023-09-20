import { useEffect } from 'react';
import axios from 'axios';

function InterceptorComponent({ setIsAuthenticated, isAuthenticated }) {

  useEffect(() => {
    // Skip interceptor if user is not authenticated
    if (!isAuthenticated) {
      return;
    }

    const interceptor = axios.interceptors.response.use(
      response => {
        return response;
      },
      error => {
        if (error.response && error.response.data.forceLogout) {
          setIsAuthenticated(false);
          // Perform client-side cleanup
          localStorage.removeItem("user");
          // Redirect to the sign-in route
          window.location.href = "/signin";
        }
        return Promise.reject(error);
      }
    );

    return () => {
      axios.interceptors.response.eject(interceptor);
    };

  }, [isAuthenticated, setIsAuthenticated]); // Re-run when isAuthenticated changes

  return null;
}

export default InterceptorComponent;


// import { useEffect } from 'react';
// import { useLocation } from 'react-router-dom';
// import axios from 'axios';

// function InterceptorComponent({ setIsAuthenticated }) {
//   const location = useLocation();

//   useEffect(() => {
//     console.log(location.pathname)
//     // Skip interceptor for specific routes
//     if (location.pathname === '/signin' || location.pathname === '/register') {
//       return;
//     }

//     const interceptor = axios.interceptors.response.use(
//       response => {
//         return response;
//       },
//       error => {
//         if (error.response && error.response.data.forceLogout) {
//           setIsAuthenticated(false);
//           // Perform client-side cleanup
//           localStorage.removeItem("user");
//           // Redirect to the sign-in route
//           window.location.href = "/signin";
//         }
//         return Promise.reject(error);
//       }
//     );

//     return () => {
//       axios.interceptors.response.eject(interceptor);
//     };

//   }, [location.pathname, setIsAuthenticated]); // Re-run when path changes

//   return null;
// }

// export default InterceptorComponent;
