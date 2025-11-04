// imports libraries(components)
import axios from "axios";
// imports hooks
import useAuth from "./useAuth";
import { useEffect } from "react";
// defines apiUrl from environment variable
const apiUrl = import.meta.env.VITE_API_BASE_URL;

// creates and exports useAxiosPrivate hook that handles private axios requests
const useAxiosPrivate = () => {
  // creates required axios auth instance with properties that requires user credentials
  const axiosAuth = axios.create({
    baseURL: apiUrl,
    withCredentials: true,
  });

  // uses custom useAuth hook to manage authentication
  const { auth, setAuth } = useAuth();
  // defines refreshing state of token
  let isRefreshing = false;
  // defines queue of failed requests
  let failedQueue = [];

  // creates helper to process queued requests after token refresh
  const processQueue = (error, response = null) => {
    failedQueue.forEach((prom) => {
      if (error) {
        prom.reject(error);
      } else {
        prom.resolve(response);
      }
    });

    failedQueue = [];
  };

  // uses useEffect hook to add interceptors to axios auth instance if authentication state changes
  useEffect(() => {
    axiosAuth.interceptors.response.use(
      (response) => response,
      async (error) => {
        console.log("⚠ Interceptor caught error:", error);
        const originalRequest = error.config;

        // if refresh token is invalid or expired shows error and rejects request using Promise
        if (
          originalRequest.url.includes("/refresh") &&
          error.response.status === 401
        ) {
          //edge case where the refresh token is invalid or expired
          console.error("❌ Refresh token has expired or is invalid.");
          return Promise.reject(error); // fail directly, no retry
        }

        if (
          error.response &&
          error.response.status === 401 &&
          !originalRequest._retry
        ) {
          // if token state is refreshing
          if (isRefreshing) {
            // retuns a promise that resolves to the original request after the token is refreshed
            return new Promise((resolve, reject) => {
              failedQueue.push({ resolve, reject });
            })
              .then(() => axiosAuth(originalRequest))
              .catch((err) => Promise.reject(err));
          }

          originalRequest._retry = true;
          isRefreshing = true;

          return new Promise((resolve, reject) => {
            axiosAuth
              .post("/refresh")
              .then(() => {
                processQueue(null);

                axiosAuth(originalRequest).then(resolve).catch(reject);
              })
              .catch((refreshError) => {
                processQueue(refreshError, null);

                localStorage.removeItem("user");
                setAuth(null); // Clear auth state
                reject(refreshError); // fail the original promise chain
              })
              .finally(() => {
                isRefreshing = false;
              });
          });
        }

        return Promise.reject(error);
      }
    );
  }, [auth]);

  return axiosAuth;
};

export default useAxiosPrivate;
