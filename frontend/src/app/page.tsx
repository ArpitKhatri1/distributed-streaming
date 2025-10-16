"use client"
import { api } from '@/lib/axios';
import axios from 'axios';
import { ChangeEvent, useState } from 'react';

export default function Home() {

  const [file, setFile] = useState<File | null>();
  const [isUploading, setIsUploading] = useState<boolean>(false);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (!selectedFile) return;
    setFile(selectedFile);
  }

  const handleFileUpload = async () => {
    if (!file) return;
    setIsUploading(true);

    const urlResponse = await api.post(`/presigned`, {
      filename: file.name,
      filetype: file.type
    })
    const awsUrl = urlResponse.data.url
    console.log(awsUrl)
    const response = await axios.put(awsUrl, file, {
      headers: {
        "Content-Type": file.type,
      }
    })

    console.log(response.data)
    setIsUploading(false);

  }




  return (
    <div>
      Upload your file here
      <br />

      <input type="file" onChange={handleFileChange} />
      <button onClick={handleFileUpload}> Upload file</button>
      <br />
      {isUploading ? ("uploading") : ("not uploading")}
    </div>
  );
}
