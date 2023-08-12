'use client';

import { Card } from '@tremor/react';
import { useState } from 'react';
import { Get, Patch } from '../../lib/http';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { useEffect } from 'react';
import {SkeletonLineItem} from "../../components/skeletons/line-items"

  export default function DetailPage({ params }: { params: { slug: string } }) {
    const [hosts, setHosts] = useState()
  
    async function GetHost() {
      const res = await Get("/hosts/" + params.slug);
      setHosts(res.response);
    }
    
    async function UpdateHost(data) {
        await Patch("/hosts/" + params.slug, data);
      }

    useEffect(() => {
      GetHost();
    }, [])
    
    const stringToIntValueConverter = (value: string) => parseInt(value, 10);
    if (!hosts) return <SkeletonLineItem/>
    return (
      <main className="p-4 md:p-10 mx-auto max-w-7xl">
        <Card className="mt-6 mb-5">
          <div>
            <Formik
              initialValues={{ name: hosts.name, address: hosts.address, port: hosts.port, username: hosts.username, password: hosts.password }}
              validate={values => {
                const errors = {};
                if (!values.name) {
                  errors.name = 'Required';
                }
                if (!values.address) {
                  errors.address = 'Required';
                }
                if (!values.port) {
                  errors.port = 'Required';
                }
                if (!values.username) {
                  errors.username = 'Required';
                }
                if (!values.password) {
                  errors.password = 'Required';
                }
                return errors;
              }}
              onSubmit={(values, { setSubmitting }) => {
                setTimeout(() => {
                  UpdateHost(JSON.stringify(values));
                  setSubmitting(false);
                }, 400);
              }}
            >
              {({ isSubmitting }) => (
                <Form>
                  <div className="mb-6">
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Hostname</label>
                    <Field type="text" name="name" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                    <ErrorMessage name="name" component="div" />
                  </div>
  
                  <div className="grid gap-6 mb-6 md:grid-cols-2">
  
                    <div>
                      <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Address</label>
                      <Field type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                      <ErrorMessage name="address" component="div" />
                    </div>
                    <div>
                      <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Port</label>
                      <Field type="number" parse={stringToIntValueConverter} name="port" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                      <ErrorMessage name="port" component="div" />
                    </div>
                    <div>
                      <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label>
                      <Field type="text" name="username" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                      <ErrorMessage name="username" component="div" />
                    </div>
                    <div>
                      <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                      <Field type="password" name="password" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                      <ErrorMessage name="password" component="div" />
                    </div>
                  </div>
  
                  <button disabled={isSubmitting} type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Save</button>
  
                </Form>
              )}
            </Formik>
          </div>
        </Card>
      </main>
    );
  }
  