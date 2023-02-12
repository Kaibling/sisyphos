import React from 'react';
import { useParams } from 'react-router-dom';
import { useState, useEffect } from 'react'
import axios from 'redaxios';
import './host.css'
import { ValOrNull,GetValue,ToObj } from './utils';
import { getToken,getURL }  from '../config';

const Host = () => {
    let { id } = useParams();
    const [apiResponse, setapiresponse] = useState({
        name: "",
        username: "",
        password: "",
        ssh_key: "",
        known_key: "",
        address: "",
        port: 0,
        tags: []
    });

    const handleScriptChange = (e, fieldname) => {
        let p = { ...apiResponse }
        p[fieldname] = e.target.value
        setapiresponse(p)
    }

    const read = name => {
        let url = getURL() + '/hosts/' + name

        axios.get(url, {
            headers: { Authorization: "Bearer " + getToken(), }
        })
            .then(response => {
                if (!response.data.success) {
                    console.log(response.data.response)
                } else {
                    setapiresponse(response.data.response);
                }
            }
            )
            .catch(error => {
                console.log("host read error ", error)
            })
    };
    const update = () => {
        apiResponse.tags =  ToObj(apiResponse.tags)
        let url = getURL() +'/hosts/' + apiResponse.name
        axios.patch(url,apiResponse, {
            headers: { Authorization: "Bearer " + getToken() }
        })
         .then(response => {
                if (!response.data.success) {
                    console.log(response.data.response)
                } else {
                    setapiresponse(response.data.response);
                }
            }
            )
            .catch(error => {
                console.log("read error ", error)
            })
    };
    useEffect(() => {
        read(id);
    }, []);


    const save = () => {
        update();
    };

    return (
        <div>
            <h3>{apiResponse.name}</h3>
            <form className='form'>
                <div className='form-group'>
                    <div className='flex-child'>
                        <label >Address</label>
                        <input type="text" className="form-control" placeholder="e.g. 10.1645.2.1" value={GetValue(apiResponse["address"])} onChange={e => handleScriptChange(e,"address")}></input>
                    </div>
                    <div className='flex-child'>
                        <label>Port</label>
                        <input type="text" className="form-control" placeholder="SSH port  e.g. 22" value={GetValue(apiResponse["port"])} onChange={e => handleScriptChange(e,"port")}></input>
                    </div>
                </div>
                <div className='form-group'>
                    <div className='flex-child'>
                        <label >Username</label>
                        <input type="text" className="form-control" placeholder="ssh username" value={GetValue(apiResponse["username"])} onChange={e => handleScriptChange(e,"username")}></input>
                    </div>
                    <div className='flex-child'>
                        <label>Password</label>
                        <input type="text" className="form-control" placeholder="ssh password"  value={GetValue(apiResponse["password"])} onChange={e => handleScriptChange(e,"password")}></input>
                    </div>
                </div >
                <div className='form-rows' >
                    <div className='single-child'>
                        <label>SSH public key</label>
                        <input type="text" className="form-control" placeholder="insert the public key for the user" value={GetValue(apiResponse["ssh_key"])} onChange={e => handleScriptChange(e,"ssh_key")}></input>
                    </div>
                    <div className='single-child'>
                        <label>Known Host key</label>
                        <input type="text" className="form-control" placeholder="known host key for ssh access" value={GetValue(apiResponse["known_key"])} onChange={e => handleScriptChange(e,"known_key")}></input>
                    </div>
                </div>


            </form>
            <button onClick={save}>save</button>
        </div>
    );
};

export default Host;