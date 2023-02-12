import { useState, useEffect } from 'react'
import './App.css'
import axios from 'redaxios';

function Overview() {
  const [tags, updateTag] = useState([])
  const [value, setValue] = useState("");
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
   const addTag = e => {      //it triggers by pressing the enter key 
    if (e.keyCode === 13) {
      if (value != "" && tags.findIndex((x) => x == value) == -1) {
        let newTags = tags.slice(); // todo copy neccessary ???
        newTags.push(value);
        updateTag(newTags);
        setValue("");
        getData();
      }
    }
  };
  const deleteTag = e => {      //it triggers by pressing the enter key 
    let newTags = tags.slice();
    let pos = newTags.findIndex((x) => x == e);
    newTags.splice(pos, 1); // TODO 1 ????
    updateTag(newTags);
  };

  const handleChange = e => {
    setValue(e.target.value);
  };

  const getData = () => {
    let url = 'http://192.168.0.94:3000/actions'
    if (tags.length > 0) {
      let filter = ""
      for (let i = 0; i < tags.length; i++) {
        filter += "tag:" + tags[i]
        if (i + 1 < tags.length) {
          filter += " "
        }
      }
      url += "?filter=" + filter
    }
    axios.get(url,{
      headers:{ Authorization: "Bearer "+ sessionStorage.getItem("s_token"),}
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
    // fetch(url)
    //   .then(function (response) {
    //     return response.json();
    //   })
    //   .then(function (myJson) {
    //     console.log(myJson)
    //      setapiresponse(myJson.response);
    //   });
  }
  useEffect(() => {
    getData();
  }, [tags]);

  return (
    <div className="App">
      <div className='tag-div'>
        <input value={value}
          onChange={handleChange}
          onKeyDown={addTag}
          placeholder="filter tags" ></input>
        <br></br>
        <table className='tag-table'><tbody>
          <tr  >
            {tags.map((value, index) => {
              return <td className='tag-card' onClick={() => { deleteTag(value) }} key={index}>{value}</td>
            })}
          </tr>
        </tbody>
        </table>
      </div>
      <table>

        <tr>
          <td>
        {apiResponse.map((value, index) => {
            return <div className='action-card' key={index}> <div><h3><a href={"/actions/"+value.name}>{value.name}</a></h3></div></div>
          })}
          </td>
        </tr>
      </table>
      {/* <div className='action-list'>

        <ul>

          {apiResponse.map((value, index) => {
            return <div className='action-card' key={index}> <div><h3><a href={"/actions/"+value.name}>{value.name}</a></h3></div></div>
          })}
        </ul>
      </div> */}
    </div>
  )
}

export default Overview


