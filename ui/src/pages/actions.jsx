import { useState, useEffect } from 'react'
import './App.css'
import axios from 'redaxios';
import { getToken,getURL }  from '../config';
import { List }  from '../components/list';



export default function Actions() {
  const [apiResponse, setapiresponse] = useState([{
    name: "who",
    script: "whoami",
    tags: [],
    triggers: [
      "ps"
    ],
    hosts: [
      {
        host_name: "ssh",
        port: "22"
      }
    ],
    variables: {},
    open: false,
  },]);

  const getData = () => {
    let url = getURL() + '/actions'
    console.log("url",url)
    axios.get(url,{
      headers:{ Authorization: "Bearer "+ getToken(),}
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
        console.log("login error ", error)
    })
  }
  useEffect(() => {
    getData();
  }, []);

  return (
    <div className="actions">
          <List data={apiResponse} collection="actions"/>
    </div>
  )
}



