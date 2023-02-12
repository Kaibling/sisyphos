export const List = (props) => {
    return (
        <ul className="list" >
            {props.data.map((v,k) => {
                let url = encodeURI("/"+props.collection+"/"+v.name).replace(/\./g, '%2E')
                return <a className="list-a" href={url}><li className="listItem" key={'"'+k+"'"} >{v.name}</li></a> })}
        </ul>
    );

};

const aaa = (data) => {
    return (
        <ul>
            {data}
        </ul>
    );

};

export default aaa

