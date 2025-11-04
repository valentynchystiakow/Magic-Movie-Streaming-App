// imports libraries(components)
import Movie from "../movie/Movie";

// creates and exports all movies component
const Movies = ({ movies, message, updateMovieReview }) => {
  return (
    // Movies wrapper block
    <div className="container mt-4">
      {/* movies row(cards) block */}
      <div className="row">
        {/* if movies array is not empty, maps each movie to Movie component, in other case displays error message */}
        {movies && movies.length > 0 ? (
          movies.map((movie) => (
            <Movie
              key={movie._id}
              movie={movie}
              updateMovieReview={updateMovieReview}
            />
          ))
        ) : (
          <h2>{message}</h2>
        )}
      </div>
    </div>
  );
};

export default Movies;
