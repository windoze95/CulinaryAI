import React from 'react';
import { useNavigate, NavLink } from 'react-router-dom';
import { Navbar, Divider } from 'react-materialize';
import './Header.css';
import LogoSvg from './logo.svg';

const Header = ({ token }) => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("user");
    navigate('/signin');
    // window.location.href = "/";
  };

  const menuItems = !token ? [
    <NavLink to="/register" key="register" className="sidenav-close">
        <i className="material-icons left">person</i>Register
    </NavLink>,
    <NavLink to="/signin" key="signin" className="sidenav-close">
      <i className="material-icons left">login</i>Sign in
    </NavLink>
  ] : [
    <NavLink to="#" key="o1" className="sidenav-close">
      User Option 1
    </NavLink>,
    <NavLink to="#" key="o2" className="sidenav-close">
      User Option 2
    </NavLink>,
    <Divider key="divider" />,
    <NavLink onClick={handleLogout} key="logout" className="sidenav-close">
      Logout
    </NavLink>
  ];

  return (
    <Navbar
      alignLinks="right"
      brand={<a className="brand-logo" href="/"><img className="logo-img" src={LogoSvg} alt="Logo" /></a>}
      id="mobile-nav"
      // menuIcon={<a href="#" data-target="mobile-nav" className="sidenav-trigger right"><i className="material-icons">menu</i></a>}
      options={{
        draggable: true,
        edge: 'right',
        inDuration: 250,
        outDuration: 200,
        preventScrolling: true,
        closeOnClick: true
      }}
    >
      {/* <Dropdown
        id="Dropdown_6"
        className="desktop-dropdown"
        options={{
          alignment: 'left',
          autoTrigger: true,
          closeOnClick: true,
          constrainWidth: true,
          container: null,
          coverTrigger: true,
          hover: false,
          inDuration: 150,
          outDuration: 250
        }}
        trigger={<a href="#!">Dropdown<i className="material-icons right">arrow_drop_down</i></a>}
      >
        {}
      </Dropdown> */}
      {menuItems}
    </Navbar>
  );
};

export default Header;
