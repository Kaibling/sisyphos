import React from 'react';
import {  Link } from "react-router-dom";
import "./nav.css";

const Navbar = (props) => {
  return (
  <div className='navigation'>
      <Link to="/"> <li className="nav-li" key="0">Sisyphos</li></Link>
      <Link to="/actions"> <li className="nav-li" key="1">Actions </li></Link>
      <Link to="/hosts"> <li className="nav-li" key="2">Hosts </li></Link>
    <li key="3" className='nav-liless'>
    <button onClick={() => props.handleLogout()}>Logout</button>
    </li>
  </div>
  
  );
}
export default Navbar;

