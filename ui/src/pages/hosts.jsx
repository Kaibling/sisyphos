import { useState, useEffect ,useContext } from 'react'
import './App.css'
import axios from 'redaxios';
import { getToken,getURL }  from '../config';
import { List }  from '../components/list';



export default function Hosts() {
  const [apiResponse, setapiresponse] = useState([
    {
      "name": "test21.local",
      "username": "test",
      "password": null,
      "ssh_key": null,
      "known_key": "",
      "address": "127.0.0.1",
      "port": 22,
      "tags": [
          "local",
          "ftp"
      ]
  },
]);

  // const handleChange = e => {
  //   setValue(e.target.value);
  // };

  const getData = () => {


    let url = getURL() + '/hosts'
    console.log("url",url)
    //if (tags.length > 0) {
      // let filter = ""
      // for (let i = 0; i < tags.length; i++) {
      //   filter += "tag:" + tags[i]
      //   if (i + 1 < tags.length) {
      //     filter += " "
      //   }
      // }
      //url += "?filter=" + filter
    //}
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
    <div className="hosts">
          <List data={apiResponse} collection="hosts"/>
    </div>
  )
}



