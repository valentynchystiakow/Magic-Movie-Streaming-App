// imports libraries(components)
import axiosClient from "../../api/axiosConfig";
import Movies from "../movies/Movies";
// imports hooks
import { useState, useEffect } from "react";

// creates and exports home component
const Home = ({ updateMovieReview }) => {
  // // hooks section
  // // uses useState hook to manage movies state
  const [movies, setMovies] = useState([]);
  // // uses useState hook to manage page loading state
  const [loading, setLoading] = useState(false);
  // // uses useState hook to manage error message state
  const [message, setMessage] = useState("");

  // uses useEffect hook to fetch movies data from database when the component mounts
  useEffect(() => {
    // defines async function to fetch movies
    const fetchMovies = async () => {
      setLoading(true);
      setMessage("");
      // try catch block to handle errors while making get request to database
      try {
        const response = await axiosClient.get("/movies");
        setMovies(response.data);
        if (response.data.length === 0) {
          setMessage("There are currently no movies available");
        }
      } catch (error) {
        console.error("Error fetching movies:", error);
        setMessage("Error fetching movies");
        // in any case sets loading state to false
      } finally {
        setLoading(false);
      }
    };
    // calls fetchMovies function
    fetchMovies();
  }, []);

  return (
    <>
      {/* Depending on loading state displays loading spinner or movies */}
      {loading ? (
        <h2> Loading ...</h2>
      ) : (
        <Movies
          movies={movies}
          message={message}
          updateMovieReview={updateMovieReview}
        />
      )}
    </>
  );
};
export default Home;
