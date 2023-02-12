import React from 'react';
import { useParams } from 'react-router-dom';
import { useState, useEffect } from 'react'
import axios from 'redaxios';
import './host.css'
import { GetValue,ToObj } from './utils';
import { getToken,getURL }  from '../config';

const Action = () => {
    let { id } = useParams();
    const [apiResponse, setapiresponse] = useState({
        name: "",
        script: "",
        tags: [],
        groups: [],
        hosts: [],
        actions: [],
        variables: {},
        fail_on_errors:false,
    });
    const handleScriptChange = (e, fieldname) => {
        let p = { ...apiResponse }
        p[fieldname] = e.target.value
        setapiresponse(p)
    }

    const read = name => {
        let url = getURL() + '/actions/' + name
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
                console.log("read error ", error)
            })
    };
    const update = () => {
        apiResponse.hosts =  ToObj(apiResponse.hosts)
        apiResponse.actions =  ToObj(apiResponse.actions)
        apiResponse.tags =  ToObj(apiResponse.tags)
        apiResponse.groups =  ToObj(apiResponse.groups)
        apiResponse.variables =  ToObj(apiResponse.variables)
        let url = getURL() +'/actions/' + apiResponse.name
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
                    <label>Tags</label>
                        <input type="text" className="form-control" placeholder="add tag" value={GetValue(apiResponse["tags"])} onChange={e => handleScriptChange(e,"tags")}></input>
                    </div>
                    <div className='flex-child'>
                    <label>Groups</label>
                        <input type="text" className="form-control" placeholder="add group" value={GetValue(apiResponse["groups"])} onChange={e => handleScriptChange(e,"groups")}></input>
                    </div>
                </div>
                <div className='form-group'>
                     <div className='flex-child'>
                        <label>Actions</label>
                            <input type="text" className="form-control" placeholder="add action" value={GetValue(apiResponse["actions"])} onChange={e => handleScriptChange(e,"actions")}></input>
                    </div>
                    <div className='flex-child'>
                        <label>Hosts</label>
                        <input type="text" className="form-control" placeholder="add host" value={GetValue(apiResponse["hosts"])} onChange={e => handleScriptChange(e,"hosts")}></input>
                    </div>
                </div>
                <div className='form-rows' >
                    <div className='single-child'>
                        <label>Script</label>
                        <textarea placeholder="add executable code" cols="7" rows="5"  value={GetValue(apiResponse["script"])} onChange={e => handleScriptChange(e,"script")} ></textarea>
                    </div>
                    <div className='single-child'>
                        <label>Variables</label>
                        <input type="text" className="form-control" placeholder="variables" value={GetValue(apiResponse["variables"])} onChange={e => handleScriptChange(e,"variables")}></input>
                    </div>
                </div>



            </form>
            <button onClick={save}>save</button>
        </div>
    );
};

export default Action;



// name: "",
// script: "",
// tags: [],
// triggers: [],
// actions: [
//     {
//         action_name: "",
//         order: "0"
//     }
// ],
// variables: {},