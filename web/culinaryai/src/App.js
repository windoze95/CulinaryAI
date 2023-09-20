import React from 'react';
import './App.css';
import { HashRouter, Route, Routes } from 'react-router-dom';
import Signin from './Signin';
import Register from './Register'; // Make sure to import your Register component
import Profile from './Profile';
import Header from './Header';

function App() {
  const token = localStorage.getItem('accessToken');

  // if (!token) {
    return (
      <HashRouter>
        <Header token={token} />
        <div className="wrapper">
          { !token ? (
            <Routes>
              <Route path="/signin" element={<Signin />} />
              <Route path="/register" element={<Register />} />
              <Route path="/*" element={<Signin />} />
            </Routes>
          ) : (
            <Routes>
              <Route path="/profile" element={<Profile />} />
              <Route path="/*" element={<Profile />} />
            </Routes>
          )}
        </div>
      </HashRouter>
  );
}

export default App;


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