import {
  Table,
  TableHead,
  TableRow,
  TableHeaderCell,
  TableBody,
  TableCell,
  Text
} from '@tremor/react';
import Link from 'next/link';


export default function HostsTable({ hosts }: any) {
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableHeaderCell>Name</TableHeaderCell>
          <TableHeaderCell>Domain</TableHeaderCell>
          <TableHeaderCell>Datacenter</TableHeaderCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {hosts?.map((host) => (
          <TableRow key={host.nameid}>
            <Link href={"/hosts/" + host.name}>
            <TableCell>{host.name}</TableCell>
            </Link>
            <TableCell>
              <Link href={"/hosts/" + host.name}>
              <Text>{host.name}</Text>
              </Link>
            </TableCell>
            <TableCell>
              <Text>{host.domain}</Text>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
