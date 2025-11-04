// imports hooks
import useAxiosPrivate from "../../hooks/useAxiosPrivate";
import { useEffect, useState } from "react";
// import libraries(components)
import Movies from "../movies/Movies";
import Spinner from "../spinner/Spinner";

// creates and exports Recommended component
const Recommended = () => {
  // hooks section
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState();
  // uses custom hook to manage private axios requests to get recommended movies
  const axiosPrivate = useAxiosPrivate();

  // uses useEffect hook to fetch recommended movies from protected recommendedmovies route when the component mounts
  useEffect(() => {
    const fetchRecommendedMovies = async () => {
      setLoading(true);
      setMessage("");

      // try catch block to handle exceptions while making private get request to database
      try {
        const response = await axiosPrivate.get("/recommendedmovies");
        // sets movies state with response data from database
        setMovies(response.data);
      } catch (error) {
        console.error("Error fetching recommended movies:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchRecommendedMovies();
  }, []);

  return (
    <>
      {/* depending on loading state returns either loading spinner or Movies component */}
      {loading ? <Spinner /> : <Movies movies={movies} message={message} />}
    </>
  );
};

export default Recommended;
