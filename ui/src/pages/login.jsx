import React, { Component } from 'react';
import axios from 'redaxios';

export default class Login extends Component{
    constructor(props) {
        super(props)
   
    this.state = {
        username: "",
        password: "",
        loginerror :""
    }
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
   
    }

    handleChange(event) {
        this.setState({
            [event.target.name]:event.target.value
        })
    }

    handleSubmit(event) {
    const {username,password} = this.state;

    axios.post("http://192.168.0.95:3000/authentication/login",
    {
        username:username,
        password:password
    },
    )
    .then(response => {
        if (!response.data.success) {   
        } else {
            this.props.handleSuccessfulAuth(response.data)
        }
    }
    )
    .catch(error => {
        console.log("login error ", error)
    })
    event.preventDefault()

}
    render() {
    return (
        <div className="loginBackground">
            Sisyphos
        <div >
            <form onSubmit={this.handleSubmit} >
                <table>
                <tbody>
                    <tr>
                        <td>
                    <input
                type="username"
                name="username"
                placeholder="username"
                value={this.state.username}
                onChange={this.handleChange}
                required
            />
            </td>
                    </tr>
                    <tr>
                    <td>
                    <input
                type="password"
                name="password"
                placeholder="password"
                value={this.state.password}
                onChange={this.handleChange}
                required
            /> </td>
                    </tr>
                    </tbody>
                </table>

            <button className="button" type="submit">Login</button>
            </form>
           
        </div>
        </div>
    );
        }
};