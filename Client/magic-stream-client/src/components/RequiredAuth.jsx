// imports libraries(components)
import { useLocation, Navigate, Outlet } from "react-router-dom";
import Spinner from "./spinner/Spinner";
// imports hooks
import useAuth from "../hooks/useAuth";

// creates and exports RequiredAuth component
const RequiredAuth = () => {
  // hooks section
  const { auth, loading } = useAuth();
  const location = useLocation();

  // if page is loading renders spinner
  if (loading) {
    return <Spinner />;
  }

  // depending on auth state returns either outlet or navigate to login page component
  return auth ? (
    <Outlet />
  ) : (
    <Navigate to="/login" state={{ from: location }} replace />
  );
};

export default RequiredAuth;
