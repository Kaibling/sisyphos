import axios, { AxiosResponse } from 'axios';
import { getSession } from "next-auth/react"
import {Toaster} from "../components/toaster.jsx";
import React from 'react';
import { render } from 'react-dom';

const backendURL = "/app";

export async function Get(url: string): Promise<AxiosResponse> {
  const session = await getSession();
  try {
    const response = await axios.get(backendURL + url, {
      headers: {
        Authorization: `Bearer ${session?.user?.accessToken}`,
      },
      withCredentials: false,
    });
    return response.data;
  } catch (error) {
    // Handle error here
    throw error;
  }
}

export async function Post(url: string, data: any): Promise<AxiosResponse> {
  const session = await getSession();
  const domNode = document.getElementById('toaster');

  try {
    const response = await axios.post(backendURL + url, data, {
      headers: {
        Authorization: `Bearer ${session?.user?.accessToken}`,
      },
    });
    if (response.data.success){
    render(<Toaster text="Created successfully" />, domNode);
    } else {
      render(<Toaster text={response.data.response} />, domNode);
    }
    return response.data;
  } catch (error) {
    // Handle error here
    throw error;
  }
}

export async function Patch(url: string, data: any): Promise<AxiosResponse> {
  const session = await getSession();
  const domNode = document.getElementById('toaster');
  try {
    const response = await axios.patch(backendURL + url, data, {
      headers: {
        Authorization: `Bearer ${session?.user?.accessToken}`,
      },
    });

   if (response.data.success){
    render(<Toaster text="Updated successfully" />, domNode);
    } else {
      render(<Toaster text={response.data.response} />, domNode);
    }
    return response.data;
  } catch (error) {
    // Handle error here
    return render(<Toaster text={error} />, domNode);
  }
}
