'use client';

import { Card } from '@tremor/react';
import { useState, useEffect } from 'react';
import HostsTable from './table';
import { Get } from '../lib/http';
import Link from 'next/link'

export default function HostsPage() {
  const [hosts, setHosts] = useState([])
  async function getHosts() {
    const res = await Get("/hosts");
    if (res.success){
    setHosts(res.response);
  }
  }
  useEffect(() => {
    getHosts();
  }, [])

  return (
    <main className="p-4 md:p-10 mx-auto max-w-7xl">
      <Link href="/hosts/new">
        <button className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
          Create new host
        </button>
      </Link>

      <Card className="mt-6 mb-5">
       {hosts.length !== 0 ? <HostsTable hosts={hosts} /> : "No hosts found.."}
      </Card>
    </main>
  );
}
