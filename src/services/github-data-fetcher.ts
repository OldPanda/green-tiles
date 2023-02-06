import type { FetchDataResponse, APIResponse, GitHubContributions } from '@/types';

const API_ENDPOINT = "https://api.old-panda.com/github/contributions?username="

export async function fetchGitHubContributions(username: string): Promise<GitHubContributions> {
  if (username.length === 0) {
    throw new Error("Username is empty.");
  }
  let url = API_ENDPOINT + username;
  let resp = await fetch(url);
  let json: APIResponse = await resp.json();
  if (resp.status !== 200) {
    throw new Error(json.error);
  }
  return <GitHubContributions>JSON.parse(json.data);
}
