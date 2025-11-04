import { useParams } from "react-router-dom";
import ReactPlayer from "react-player";
import "./StreamMovie.css";

// creates and exports StreamMovie function
const StreamMovie = () => {
  // hooks section
  // uses useParams hook to get yt_id from url
  let params = useParams();
  let key = params.yt_id;

  return (
    <div className="react-player-container">
      {/* if key is not null, renders react player component */}
      {key != null ? (
        <ReactPlayer
          controls="true"
          playing={true}
          url={`https://www.youtube.com/watch?v=${key}`}
          width="100%"
          height="100%"
        />
      ) : null}
    </div>
  );
};

export default StreamMovie;
