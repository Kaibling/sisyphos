import { Component } from 'react';
import {Routes, BrowserRouter , Route  } from 'react-router-dom';
import Navbar from './components/nav';
import PageNotFound from './pages/pagenotfound';
import Login from './pages/login';
import Hosts from './pages/hosts';
import Host from './pages/host';
import Action from './pages/action';
import Overview from './pages/overview';
import Actions from './pages/actions';
import InitStorage from './config';

export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
       loggedIn: false,
       user:""
       }
       InitStorage();
       this.handleSuccessfulAuth = this.handleSuccessfulAuth.bind(this);
       this.handleLogout = this.handleLogout.bind(this);
      }
      checkLoginStatus() {
        var userToken = sessionStorage.getItem("s_token");
        var userName = sessionStorage.getItem("s_user");
      if (userToken != null) {
        this.setState({
          loggedIn:true,
          user:userName
        })
      }
      }
      componentDidMount() {
        this.checkLoginStatus();
      }
      
      handleLogin(data) {
        this.setState({
          loggedIn:true,
          user:data.response.name
        })
        sessionStorage.setItem("s_token", data.response.token[0]);
        sessionStorage.setItem("s_user", data.response.name);
      }
      
      handleSuccessfulAuth(data) {
        this.handleLogin(data)      
      }
      handleLogout(){
        sessionStorage.removeItem("s_token")
        sessionStorage.removeItem("s_user")
        this.setState({
          loggedIn:false,
          user:""
        })
      }
      
  render() {
    return (
      
      <div >
         {this.state.loggedIn ? (
      <BrowserRouter >
        <Navbar handleLogout={this.handleLogout} />
        <Routes >
          <Route path='/' exact element={<Overview />} />
          <Route exact path='/hosts' element={<Hosts/>} />
          <Route exact path='/actions' element={<Actions/>} />
          <Route path='/hosts/:id' element={<Host/>} />
          <Route path='/actions/:id' element={<Action/>} />
          <Route path='*' element={<PageNotFound/>} />
        </Routes>
    </BrowserRouter>
       ) : (     
        <Login handleSuccessfulAuth={this.handleSuccessfulAuth}  /> 
        )}
    </div>
  
    )
  }
}
