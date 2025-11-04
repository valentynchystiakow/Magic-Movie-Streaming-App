// imports libraries(components)
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import logo from "../../assets/MagicStreamLogo.png";
// imports hooks
import { useNavigate, NavLink, Link } from "react-router-dom";
import useAuth from "../../hooks/useAuth";

// creates and exports Header component
const Header = ({ handleLogout }) => {
  // hooks section
  // uses useNavigate hook to navigate between pages
  const navigate = useNavigate();
  const { auth } = useAuth();

  return (
    // header wrapper block
    <Navbar
      bg="dark"
      variant="dark"
      expand="lg"
      stick="top"
      className="shadow-sm"
    >
      {/* Brand block */}
      <Container>
        <Navbar.Brand>
          <img
            alt=""
            src={logo}
            width="30"
            height="30"
            className="d-inline-block align-top me-2"
          />
          Magic Stream
        </Navbar.Brand>

        {/* Navigation block */}
        <Navbar.Toggle aria-controls="main-navbar-nav" />
        <Navbar.Collapse>
          <Nav className="me-auto">
            {/* Home link */}
            <Nav.Link as={NavLink} to="/">
              Home
            </Nav.Link>
            {/* Recommended Link */}
            <Nav.Link as={NavLink} to="/recommended">
              Recommended
            </Nav.Link>
          </Nav>

          {/* depending on auth state returns either logout button or login and register buttons */}
          <Nav className="ms-auto align-items-center">
            {auth ? (
              <>
                {/* Welcome message */}
                <span className="me-3 text-light">
                  Hello, <strong>{auth.first_name}</strong>
                </span>
                {/* Logout button */}
                <Button
                  variant="outline-light"
                  size="sm"
                  onClick={handleLogout}
                >
                  Logout
                </Button>
              </>
            ) : (
              <>
                {/* Login and register buttons */}
                <Button
                  variant="outline-info"
                  size="sm"
                  className="me-2"
                  onClick={() => navigate("/login")}
                >
                  Login
                </Button>
                <Button
                  variant="info"
                  size="sm"
                  onClick={() => navigate("/register")}
                >
                  Register
                </Button>
              </>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};
export default Header;
