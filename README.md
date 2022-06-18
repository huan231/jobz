# jobz

_Monitoring system for k8s cron jobs_

## Installation

[Helm](https://helm.sh) must be installed to use the charts. Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

```sh
helm repo add jobz https://huan231.github.io/jobz/
```

And install the jobz chart:

```sh
helm install jobz/jobz --generate-name
```

## License

This project is licensed under the terms of the [MIT license](https://github.com/huan231/jobz/blob/master/LICENSE).
