import React, { createContext, useContext, useState, useEffect } from 'react';
import './App.css';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Signin from './Signin';
import Register from './Register';
import Profile from './Profile';
import Header from './Header';
import axios from 'axios';

const AuthContext = createContext();

axios.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response && error.response.data.forceLogout) {
      // Perform client-side cleanup
      localStorage.removeItem("user");
      // Redirect to the sign-in route
      window.location.href = "/signin";
    }
    return Promise.reject(error);
  }
);

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isVerifying, setIsVerifying] = useState(true);

  useEffect(() => {
    // Verify the JWT token in the HTTP-only cookie
    axios.get('/api/v1/users/verify', { withCredentials: true })
      .then(response => {
        if (response.data.isAuthenticated) {
          setIsAuthenticated(true);
        }
      })
      .catch(error => {
        setIsAuthenticated(false);
      })
      .finally(() => {
        setIsVerifying(false); // Set to false once verification is done
      });
  }, []); // Empty dependency array means this useEffect runs once when the component mounts

  return (
    <AuthContext.Provider value={{ isAuthenticated }}>
      <BrowserRouter>
        <Header />
        <div className="wrapper">
          <Routes>
            {!isAuthenticated ? (
              <>
                <Route path="/signin" element={<Signin />} />
                <Route path="/register" element={<Register />} />
                <Route path="/*" element={<Signin />} />
              </>
            ) : (
              <>
                <Route path="/profile" element={<Profile />} />
                <Route path="/*" element={<Profile />} />
              </>
            )}
          </Routes>
        </div>
      </BrowserRouter>
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}

export default App;

// const AuthContext = createContext();

// function App() {
//   const token = localStorage.getItem('accessToken');

//   return (
//     <AuthContext.Provider value={{ token }}>
//       <BrowserRouter>
//         <Header />
//         <div className="wrapper">
//           <Routes>
//             {!token ? (
//               <>
//                 <Route path="/signin" element={<Signin />} />
//                 <Route path="/register" element={<Register />} />
//                 <Route path="/*" element={<Signin />} />
//               </>
//             ) : (
//               <>
//                 <Route path="/profile" element={<Profile />} />
//                 <Route path="/*" element={<Profile />} />
//               </>
//             )}
//           </Routes>
//         </div>
//       </BrowserRouter>
//     </AuthContext.Provider>
//   );
// }

// export function useAuth() {
//   return useContext(AuthContext);
// }

// export default App;

// import React from 'react';
// import './App.css';
// import { HashRouter, Route, Routes } from 'react-router-dom';
// import Signin from './Signin';
// import Register from './Register'; // Make sure to import your Register component
// import Profile from './Profile';
// import Header from './Header';

// function App() {
//   const token = localStorage.getItem('accessToken');

//   // if (!token) {
//     return (
//       <HashRouter>
//         <Header token={token} />
//         <div className="wrapper">
//           { !token ? (
//             <Routes>
//               <Route path="/signin" element={<Signin />} />
//               <Route path="/register" element={<Register />} />
//               <Route path="/*" element={<Signin />} />
//             </Routes>
//           ) : (
//             <Routes>
//               <Route path="/profile" element={<Profile />} />
//               <Route path="/*" element={<Profile />} />
//             </Routes>
//           )}
//         </div>
//       </HashRouter>
//   );
// }

// export default App;


// import React from 'react';
// import './App.css';
// import { BrowserRouter, Route, Routes } from 'react-router-dom';
// import Signin from './Signin';
// import Profile from './Profile';

// function App() {
//   const token = localStorage.getItem('accessToken');

//   if(!token) {
//     return <Signin />
//   }

//   return (
//     <div className="wrapper">
//       <BrowserRouter>
//         <Routes>
//           <Route path="/profile">
//             <Profile />
//           </Route>
//           <Route path="/">
//             <Profile />
//           </Route>
//         </Routes>
//       </BrowserRouter>
//     </div>
//   );
// }

// export default App;