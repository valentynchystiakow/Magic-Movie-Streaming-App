// imports hook
import { useContext } from "react";
// imports context
import AuthContext from "../context/AuthProvider.jsx";

// creates and export useAuth hook(custom hook)
const useAuth = () => {
  return useContext(AuthContext);
};

export default useAuth;
