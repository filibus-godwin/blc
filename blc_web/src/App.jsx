import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { getFirestore, collection, getDocs, onSnapshot } from 'firebase/firestore'
import { initializeApp } from "firebase/app";

const firebaseConfig = {
  apiKey: "AIzaSyCZ4NbbA4y-F-ZbpJeU-h-7mpfohWpA18g",
  authDomain: "voting-project-46c2f.firebaseapp.com",
  databaseURL: "https://voting-project-46c2f-default-rtdb.firebaseio.com",
  projectId: "voting-project-46c2f",
  storageBucket: "voting-project-46c2f.appspot.com",
  messagingSenderId: "223278576698",
  appId: "1:223278576698:web:959fee6e1ccd6fc3f03103"
};

// Initialize Firebase
export const app = initializeApp(firebaseConfig);

const db = getFirestore(app)

function App() {
  const [voteData, setVoteData] = useState({});
  const candidateData = {
    1: {name:"Asiwaju Bola Tinibu", party:'APC'},
    2: {name:"Atiku Abubakar", party : "PDP"},
    3: {name:"Peter Gregory Obi", party: 'LP'},
  }

  useEffect(() => {
    onSnapshot(collection(db, "votes"), (snapshot) => {
      const data = {
        1: 0,
        2: 0,
        3: 0,
      };
      snapshot.forEach((doc) => {
        const partyId = doc.data().partyId;
        data[partyId] = data[partyId] + 1
      })
      setVoteData(data);
    })
  }, [])

  return (
    <div className="App">
      <div className="read-the-docs">
        <p style={{fontWeight:'bolder', fontSize:'4em', margin:0}}>Election Results</p>
      </div>
      <div className="read-the-docs" style={{backgroundColor: '#0e0e0e', marginBottom:2}}>
        <p style={{ width: 120, display: 'inline-block', }}>Party</p>
        <p style={{ width: 400, display: 'inline-block',  }}>Name</p>
        <p style={{ width: 100, display: 'inline-block',  }}>Votes</p>
      </div>
      <div className="read-the-docs">
        {
          Object.keys(candidateData).map((d) => {
            return <div key={d} style={{ backgroundColor: '#0e0e0e' }}>
              <p style={{ width: 120, display: 'inline-block', backgroundColor: '#0e0e0e',}}>{candidateData[d].party}</p>
              <p style={{ width: 400, display: 'inline-block', backgroundColor: '#0e0e0e', }}>{candidateData[d].name}</p>
              <p style={{ width: 100, display: 'inline-block', backgroundColor: '#0e0e0e', }}>{voteData[d] ?? 0}</p>
            </div>
          })
        }
      </div>
    </div>
  )
}

export default App
