// imports libraries(components)
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosClient from "../../api/axiosConfig";
import logo from "../../assets/MagicStreamLogo.png";
// imports hooks
import { useState } from "react";
import { useNavigate, Link, useLocation } from "react-router-dom";
import useAuth from "../../hooks/useAuth";

// creates and exports Login component
const Login = () => {
  // hooks section
  // uses custom hook to manage authentication
  const { setAuth } = useAuth("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  // uses location hook to access the location object
  const location = useLocation();
  // uses navigate hook to navigate to another route(page)
  const navigate = useNavigate();

  // defines from variable - used to redirect user to previous page(to home page by default)
  const from = location.state?.from?.pathname || "/";

  // creates async function that handles login form submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // try catch block to handle exceptions while making post request to login user
    try {
      const response = await axiosClient.post("/login", { email, password });
      console.log(response.data);
      if (response.data.error) {
        setError(response.data.error);
        return;
      }
      // console.log(response.data);
      // uses setAuth hook to update authentication state with user data
      setAuth(response.data);

      // sets updated user data in local storage
      localStorage.setItem("user", JSON.stringify(response.data));
      // Handles successful login (e.g., store token, redirect)
      navigate(from, { replace: true });
      // navigate('/');
    } catch (err) {
      console.error(err);
      setError("Invalid email or password");
    } finally {
      setLoading(false);
    }
  };
  return (
    <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
      {/* Login card section */}
      <div
        className="login-card shadow p-4 rounded bg-white"
        style={{ maxWidth: 400, width: "100%" }}
      >
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="fw-bold">Sign In</h2>
          <p className="text-muted">
            Welcome back! Please login to your account.
          </p>
        </div>
        {error && <div className="alert alert-danger py-2">{error}</div>}
        {/* Login form section */}
        <Form onSubmit={handleSubmit}>
          {/* Login form input fields */}
          <Form.Group controlId="formBasicEmail" className="mb-3">
            {/* email field */}
            <Form.Label>Email address</Form.Label>
            <Form.Control
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoFocus
            />
          </Form.Group>

          <Form.Group controlId="formBasicPassword" className="mb-3">
            {/* password field */}
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </Form.Group>
          {/* login button */}
          <Button
            variant="primary"
            type="submit"
            className="w-100 mb-2"
            disabled={loading}
            style={{ fontWeight: 600, letterSpacing: 1 }}
          >
            {/* displays loading spinner or login text based on loading state */}
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                ></span>
                Logging in...
              </>
            ) : (
              "Login"
            )}
          </Button>
        </Form>
        {/* register link info */}
        <div className="text-center mt-3">
          <span className="text-muted">Don't have an account? </span>
          <Link to="/register" className="fw-semibold">
            Register here
          </Link>
        </div>
      </div>
    </Container>
  );
};
export default Login;
