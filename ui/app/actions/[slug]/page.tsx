'use client';

import { Card } from '@tremor/react';
import { useState } from 'react';
import { Get, Patch } from '../../lib/http';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { useEffect } from 'react';
import { SkeletonLineItem } from "../../components/skeletons/line-items"
import { Monaco } from "../../components/monaco.jsx";
import { MultiSelect, SelectCreatable } from "../../components/components.jsx";

export default function DetailPage({ params }: { params: { slug: string } }) {
  const [hosts, setHosts] = useState()
  const [host, setHost] = useState()
  const [tags, setTags] = useState()

  async function GetHosts() {
    const res = await Get("/hosts/");
    setHosts(res.response.map(n=>({value: n.name, label: n.name})));
  }

  async function GetHost() {
    const res = await Get("/actions/" + params.slug);
    setHost(res.response);
  }
  async function GetTags() {
    const res = await Get("/tags/");
    console.log(res.response)
    setTags(res.response.map(n=>({value: n.name, label: n.name})));
  }
  async function UpdateHost(data) {
    await Patch("/actions/" + params.slug, data);
  }
  
  const handleScriptChange = (data, setFieldValue) => {
    setFieldValue('script', data);
  };

  const handleHostChange = (data, setFieldValue) => {
    setFieldValue('hosts', data.map((n, index) => ({
      name: n.value,
      order: index + 1
    })));
  };

  const handleTagsChange = (data, setFieldValue) => {
    setFieldValue('tags', data.map(n=>(n.value)));
  };

  useEffect(() => {
    GetHost();
    GetHosts();
    GetTags();
  }, [])

  const stringToIntValueConverter = (value: string) => parseInt(value, 10);
  if (!host && !hosts) return <SkeletonLineItem />
  return (
    <main className="p-4 md:p-10 mx-auto max-w-7xl">
      <Card className="mt-6 mb-5">
        <div>
          <Formik
            initialValues={{ name: host.name, script: host.script, port: host.port, username: host.username, password: host.password }}
            validate={values => {
              const errors = {};
              if (!values.name) {
                errors.name = 'Required';
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
                  <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Name</label>
                  <Field type="text" name="name" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                  <ErrorMessage name="name" component="div" />
                </div>

                <div className="grid gap-6 mb-6 md:grid-cols-1">
                  <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Script</label>
                    <Monaco
                      onChange={value => handleScriptChange(value, setFieldValue)}
                      value={host.script}
                    />
                    <ErrorMessage name="script" component="div" />
                  </div>
                </div>
                <div className="grid gap-6 mb-6 md:grid-cols-2">
                  <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Hosts</label>
                    <MultiSelect
                      onChange={value => handleHostChange(value, setFieldValue)}
                      options={hosts}
                      className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    />
                    <ErrorMessage name="hosts" component="div" />
                  </div>
                  <div>
                    <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Tags</label>
                    <SelectCreatable
                      onChange={value => handleTagsChange(value, setFieldValue)}
                      options={tags}
                      className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    />
                    <ErrorMessage name="tags" component="div" />
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
