import { pb } from "$lib/pocketbase";
import { createMutation } from "@tanstack/svelte-query";
import type { ClientResponseError } from "pocketbase";
import toast from "svelte-french-toast";
import { push } from "svelte-spa-router";

export interface LoginData {
  email: string;
  password: string;
}

const login = async (data: LoginData): Promise<void> => {
  await pb.collection("users").authWithPassword(data.email, data.password);
};

export const useLogin = () =>
  createMutation<void, ClientResponseError, LoginData>({
    mutationFn: login,
    onSuccess: () => {
      push("/");
      toast.success("Logged In!");
    },
    onError: (err) => {
      toast.error(err.message);
    },
  });

export interface SignupData {
  email: string;
  name: string;
  password: string;
  password_confirm: string;
}

const signup = async (data: SignupData) => {
  const request = {
    email: data.email,
    emailVisibility: false,
    password: data.password,
    passwordConfirm: data.password_confirm,
    name: data.name,
    role: "student",
  };

  await pb.collection("users").create(request);
};

export const useSignup = () =>
  createMutation<void, ClientResponseError, SignupData>({
    mutationFn: signup,
    onSuccess: () => {
      push("/");
      toast.success("Sign Up Sucessful!");
    },
    onError: (err) => {
      console.log(err.data);
      toast.error(err.message);
    },
  });
