// imports libraries(components)
import Button from "react-bootstrap/Button";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCirclePlay } from "@fortawesome/free-solid-svg-icons";
// imports styles
import "./Movie.css";

// creates and exports reusable movie component
const Movie = ({ movie, updateMovieReview }) => {
  return (
    // movie wrapper block
    <div className="col-md-4 mb-4" key={movie._id}>
      {/* Link block to navigate to stream movie page */}
      <Link
        to={`/stream/${movie.youtube_id}`}
        style={{ textDecoration: "none", color: "inherit" }}
      >
        {/* Movie card block */}
        <div className="card h-100 shadow-sm movie-card">
          {/* Poster image block */}
          <div style={{ position: "relative" }}>
            <img
              src={movie.poster_path}
              className="card-img-top"
              style={{ objectFit: "contain", height: "250px", width: "100%" }}
              alt={movie.title}
            />
            {/* Play icon block */}
            <span className="play-icon-overlay">
              <FontAwesomeIcon icon={faCirclePlay} />
            </span>
          </div>
          {/*Сard body block  */}
          <div className="card-body d-flex flex-column">
            {/* Card title block */}
            <h5 className="card-title">{movie.title}</h5>
            {/* Movie imdb id block */}
            <p className="card-text mb-2">{movie.imdb_id}</p>
          </div>
          {/* displayes movie ranking if it exists */}
          {movie.ranking?.ranking_name && (
            <span
              className="badge bg-dark m-3 p-2"
              style={{ fontSize: "1rem" }}
            >
              {movie.ranking.ranking_name}
            </span>
          )}
          {/* if movie review exists, displays update review button, in other case displays create review button */}
          {updateMovieReview && (
            <Button
              variant="outline-info"
              onClick={(e) => {
                e.preventDefault();
                updateMovieReview(movie.imdb_id);
              }}
              className="m-3"
            >
              Review
            </Button>
          )}
        </div>
      </Link>
    </div>
  );
};

export default Movie;
