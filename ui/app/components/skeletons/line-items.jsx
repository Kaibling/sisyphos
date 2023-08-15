import { Card } from '@tremor/react';

export const SkeletonLineItem = () => {
    return (
        <main className="p-4 md:p-10 mx-auto max-w-7xl">
            <Card className="mt-6 mb-5">
                <div className="animate-pulse">
                    <div>
                        <div>
                            <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"><div className="h-3 w-[120px] bg-gray-200" /></label>
                            <div className='relative'>
                                <input type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                                <div className="absolute h-3 w-[94%] bg-gray-200 inside-input" />
                            </div>
                        </div>
                        <div className="grid gap-6 mb-6 md:grid-cols-2 mt-5">
                        <div>
                            <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"><div className="h-3 w-[120px] bg-gray-200" /></label>
                            <div className='relative'>
                                <input type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                                <div className="absolute h-3 w-[94%] bg-gray-200 inside-input" />
                            </div>
                        </div>
                        <div>
                            <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"><div className="h-3 w-[120px] bg-gray-200" /></label>
                            <div className='relative'>
                                <input type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                                <div className="absolute h-3 w-[94%] bg-gray-200 inside-input" />
                            </div>
                        </div>
                        </div>
                        <div className="grid gap-6 mb-6 md:grid-cols-2 mt-5">
                        <div>
                            <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"><div className="h-3 w-[120px] bg-gray-200" /></label>
                            <div className='relative'>
                                <input type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                                <div className="absolute h-3 w-[94%] bg-gray-200 inside-input" />
                            </div>
                        </div>
                        <div>
                            <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"><div className="h-3 w-[120px] bg-gray-200" /></label>
                            <div className='relative'>
                                <input type="text" name="address" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" />
                                <div className="absolute h-3 w-[94%] bg-gray-200 inside-input" />
                            </div>
                        </div>
                        </div>
                        <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"><div className="h-5 w-[50px]" /></button>

                    </div>
                </div>
            </Card>
        </main>
    )
}

export default SkeletonLineItem
