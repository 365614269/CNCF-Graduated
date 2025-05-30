import {Artifact, ArtifactRepository, NodeStatus, WorkflowStatus} from './models';

export const nodeArtifacts = (node: NodeStatus, ar: ArtifactRepository) =>
    (node.inputs?.artifacts || [])
        .map(a => ({
            ...a,
            artifactNameDiscriminator: 'input'
        }))
        .concat((node.outputs?.artifacts || []).map(a => ({...a, artifactNameDiscriminator: 'output'})))
        .map(a => ({
            ...a,
            urn: artifactURN(a, ar),
            key: artifactKey(a),
            nodeId: node.id
        }))
        .map(a => ({
            ...a,
            // trim trailing slash to get the correct filename for a directory
            filename: a.key.replace(/\/$/, '').split('/').pop()
        }));

export function artifactURN<A extends Artifact>(a: A, ar: ArtifactRepository) {
    if (a.gcs) {
        return 'artifact:gcs:' + (a.gcs.endpoint || ar?.gcs?.endpoint) + ':' + (a.gcs.bucket || ar?.gcs?.bucket) + ':' + a.gcs.key;
    } else if (a.git) {
        return 'artifact:git:' + a.git.repo + ':' + (a.git.revision || a.git.branch || 'HEAD');
    } else if (a.http) {
        return 'artifact:http::' + a.http.url;
    } else if (a.s3) {
        return 'artifact:s3:' + (a.s3.endpoint || ar?.s3?.endpoint) + ':' + (a.s3.bucket || ar?.s3?.bucket) + ':' + a.s3.key;
    } else if (a.oss) {
        return 'artifact:oss:' + (a.oss.endpoint || ar?.oss?.endpoint) + ':' + (a.oss.bucket || ar?.oss?.bucket) + ':' + a.oss.key;
    } else if (a.raw) {
        return 'artifact:raw:' + a.raw.data;
    } else if (a.azure) {
        return 'artifact:azure:' + (a.azure.endpoint || ar?.azure?.endpoint) + ':' + (a.azure.container || ar?.azure?.container) + ':' + a.azure.blob;
    }
    return 'artifact:unknown';
}

export function artifactRepoHasLocation(ar: ArtifactRepository) {
    if (ar.gcs) {
        return ar.gcs.bucket !== '' && ar.gcs.key !== '';
    } else if (ar.git) {
        return ar.git.repo !== '';
    } else if (ar.http) {
        return ar.http.url !== '';
    } else if (ar.s3) {
        return ar.s3.endpoint !== '' && ar.s3.bucket !== '' && ar.s3.key !== '';
    } else if (ar.oss) {
        return ar.oss.bucket !== '' && ar.oss.endpoint !== '' && ar.oss.key !== '';
    } else if (ar.raw) {
        return true;
    } else if (ar.azure) {
        return ar.azure.container !== '' && ar.azure.blob !== '';
    }
}

export function artifactKey<A extends Artifact>(a: A) {
    if (a.gcs?.key) {
        return a.gcs.key;
    } else if (a.git?.repo) {
        return a.git.repo + '#' + (a.git.revision || 'HEAD');
    } else if (a.http?.url) {
        return a.http.url;
    } else if (a.s3?.key) {
        return a.s3.key;
    } else if (a.oss?.key) {
        return a.oss.key;
    } else if (a.raw) {
        return 'raw';
    } else if (a.azure?.blob) {
        return a.azure.blob;
    }
    return 'unknown';
}

export function findArtifact(status: WorkflowStatus, urn: string) {
    const artifacts: (Artifact & {nodeId: string; artifactNameDiscriminator: string})[] = [];

    Object.values(status.nodes || {}).map(node => {
        return nodeArtifacts(node, status.artifactRepositoryRef?.artifactRepository)
            .filter(ad => ad.urn === urn)
            .forEach(ad => artifacts.push(ad));
    });

    return artifacts.length >= 0 && artifacts[artifacts.length - 1];
}
