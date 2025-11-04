// imports libraries(components)
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosClient from "../../api/axiosConfig";
// imports hooks
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { useEffect } from "react";
// imports assets
import logo from "../../assets/MagicStreamLogo.png";

// creates and exports Register component
const Register = () => {
  // hooks section
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [favouriteGenres, setFavouriteGenres] = useState([]);
  const [genres, setGenres] = useState([]);

  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // creates function that handles genre change
  const handleGenreChange = (e) => {
    // defines genre options
    const options = Array.from(e.target.selectedOptions);
    // changes state of favourite genres based on genre options
    setFavouriteGenres(
      options.map((opt) => ({
        genre_id: Number(opt.value),
        genre_name: opt.label,
      }))
    );
  };

  // creates async function that handles registration form submission
  const handleSubmit = async (e) => {
    // prevents default form submission
    e.preventDefault();
    setError(null);
    // defines default user role
    const defaultRole = "USER";

    console.log(defaultRole);

    // checks if passwords match
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setLoading(true);

    // uses try catch block ho handle erros while making post request to register user in database
    try {
      // defines payload from form
      const payload = {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        role: defaultRole,
        favourite_genres: favouriteGenres,
      };
      // uses axiosClient to make post request with payload data
      const response = await axiosClient.post("/register", payload);
      if (response.data.error) {
        setError(response.data.error);
        return;
      }
      // Registration successful, redirects to login
      navigate("/login", { replace: true });
    } catch (err) {
      setError("Registration failed. Please try again.");
      // in any case sets loading state to false
    } finally {
      setLoading(false);
    }
  };
  // uses useEffect hooks to get genres when component mounts
  useEffect(() => {
    // defines async function to fetch genres
    const fetchGenres = async () => {
      // try catch block to handle exceptions while making get request to database
      try {
        const response = await axiosClient.get("/genres");
        setGenres(response.data);
      } catch (error) {
        console.error("Error fetching genres:", error);
      }
    };
    // calls fetchGenres function
    fetchGenres();
  }, []);

  return (
    <Container className="login-container d-flex align-items justify-content-center min-vh-100">
      {/* Register card section */}
      <div
        className="login-card shadow p-4 rounded bg-white"
        style={{ maxWidth: 400, width: "100%" }}
      >
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="fw-bold">Register</h2>
          <p className="text-muted">Create your Magic Movie Stream account.</p>
          {error && <div className="alert alert-danger py-2">{error}</div>}
        </div>
        {/* Submit form section */}
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3">
            {/* first name input */}
            <Form.Label>First Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter first name"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-3">
            {/* last name input */}
            <Form.Label>Last Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter last name"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-3">
            {/* email input  */}
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-3">
            {/* password input */}
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </Form.Group>
          <Form.Group className="mb-3">
            {/* confirm password input */}
            <Form.Label>Confirm Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Confirm Password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              isInvalid={!!confirmPassword && password !== confirmPassword}
            />
            {/* error message if passwords do not match */}
            <Form.Control.Feedback type="invalid">
              Passwords do not match.
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group>
            {/* Select genres form */}
            <Form.Select
              multiple
              value={favouriteGenres.map((g) => String(g.genre_id))}
              onChange={handleGenreChange}
            >
              {/* maps through genres and creates option for each */}
              {genres.map((genre) => (
                <option
                  key={genre.genre_id}
                  value={genre.genre_id}
                  label={genre.genre_name}
                >
                  {genre.genre_name}
                </option>
              ))}
            </Form.Select>
            {/* help text */}
            <Form.Text className="text-muted">
              Hold Ctrl (Windows) or Cmd (Mac) to select multiple genres.
            </Form.Text>
          </Form.Group>
          {/* register button */}
          <Button
            variant="primary"
            type="submit"
            className="w-100 mb-2"
            disabled={loading}
            style={{ fontWeight: 600, letterSpacing: 1 }}
          >
            {/* displays loading spinner while registering */}
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                Registering...
              </>
            ) : (
              "Register"
            )}
          </Button>
        </Form>
      </div>
    </Container>
  );
};

export default Register;
