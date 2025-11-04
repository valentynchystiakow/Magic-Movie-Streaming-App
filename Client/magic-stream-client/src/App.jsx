// imports libraries(components)
import "./App.css";
import Home from "./components/home/Home";
import Header from "./components/header/Header";
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import RequiredAuth from "./components/RequiredAuth";
import Layout from "./components/Layout";
import Recommended from "./components/recommended/Recommended";
import axiosClient from "./api/axiosConfig";
import Review from "./components/review/Review";
import { Route, Routes } from "react-router-dom";
import StreamMovie from "./components/stream/StreamMovie";
// imports hooks
import useAuth from "./hooks/useAuth";
import { useNavigate } from "react-router-dom";

// creates and exports App function that renders all app components
function App() {
  // uses custom useAuth hook to manage authentication
  const { auth, setAuth } = useAuth("");
  // uses usenNavigate hook to navigate between pages
  const navigate = useNavigate();

  // creates async function that handles logout
  const handleLogout = async () => {
    // try catch block to handle exceptions while making post request to logout user
    try {
      const response = await axiosClient.post("/logout", {
        user_id: auth.user_id,
      });
      console.log(response.data);
      // sets auth state to null
      setAuth(null);
      // localStorage.removeItem('user');
      console.log("User logged out");
    } catch (error) {
      console.error("Error logging out:", error);
    }
  };

  // creates function that updates movie review and navigates to review page
  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };

  return (
    <>
      {/* Header component */}
      <Header handleLogout={handleLogout} />
      {/* Wraps components in routes for routing between pages */}
      <Routes path="/" element={<Layout />}>
        <Route
          path="/"
          element={<Home updateMovieReview={updateMovieReview} />}
        ></Route>
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        {/* Wraps protected routes with RequiredAuth component  */}
        <Route element={<RequiredAuth />}>
          <Route path="/recommended" element={<Recommended />}></Route>
          <Route path="/review/:imdb_id" element={<Review />}></Route>
          <Route path="/stream/:yt_id" element={<StreamMovie />}></Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
