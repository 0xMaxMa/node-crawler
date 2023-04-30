import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import "../App.css";

interface NodeData {
  ID: string;
  IP: string;
  ConnType: string;
  OS: string;
  Client: string;
  Language: string;
  Validator: boolean;
}

const initialData: NodeData = {
  ID: '',
  IP: '',
  ConnType: '',
  OS: '',
  Client: '',
  Language: '',
  Validator: false,
};

function ViewAll() {

  let [currentPage, setCurrentPage] = useState(0);

  const [itemsPerPage, setItemsPerPage] = useState(10);

  let [data, setData] = useState<NodeData[]>([initialData]);

  let [slicedData, setSlicedData] = useState<NodeData[]>([initialData]);

  var renderPageNumbers: any;

  useEffect(() => {
    async function getAllData() {
      const response = await fetch(`v1/all`);
      const json: NodeData[] = await response.json();
      setData(json);
    }

    getAllData();
  }, []);

  useEffect(() => {
    // set page numbers for table
    let pageNumbers: any[] = [];
    for (let i = 1; i <= Math.ceil(data.length / itemsPerPage); i++) {
      pageNumbers.push(i);
    }
    renderPageNumbers = pageNumbers.map(number => {
      return (
        <li key={number}>
          <button onClick={() => setCurrentPage(number)}>
            {number}
          </button>
        </li>
      );
    });

    setCurrentPage(1);
  }, [data, itemsPerPage]);

  useEffect(() => {
    const indexOfLastItem = currentPage * itemsPerPage;
    const indexOfFirstItem = indexOfLastItem - itemsPerPage;
    const currentItems = data.slice(indexOfFirstItem, indexOfLastItem);

    setSlicedData(currentItems);
    
  }, [data, currentPage, itemsPerPage]);

  return (
    <div className="container">
      {currentPage > 1 && <button onClick={() => setCurrentPage(currentPage - 1)}>Previous</button>}
      
      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>IP</th>
            <th>Connection Type</th>
            <th>Operating System</th>
            <th>Client</th>
            <th>Language</th>
            <th>Is Validator</th>
          </tr>
        </thead>
        <tbody>
          {slicedData.map(item => (
            <tr key={item.ID}>
              <td>{item.ID}</td>
              <td>{item.IP}</td>
              <td>{item.ConnType}</td>
              <td>{item.OS}</td>
              <td>{item.Client}</td>
              <td>{item.Language}</td>
              <td>{item.Validator ? "Yes" : "No"}</td>
            </tr>
          ))}
        </tbody>
      </table>
      
      <button className='float-right' onClick={() => setCurrentPage(currentPage + 1)}>Next</button>

      <Link to="/">
        <button> View Chart</button>
      </Link>
    </div>
  );
}

export default ViewAll;