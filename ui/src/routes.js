import { mount, route, redirect, map, compose, withView, withData } from 'navi';

import api from './api';
import { useEffect } from 'react';
import { useNavigation } from 'react-navi';

export default mount({
  '/': redirect('/projects'),

  '/signup': route({
    title: 'Sign Up',
    getView: () => import('./containers/signup'),
  }),
  '/login': map(async (request, context) =>
    context.currentUser
      ? redirect(
          request.params.redirectTo
            ? decodeURIComponent(request.params.redirectTo)
            : '/projects'
        )
      : route({
          title: 'Log In',
          getData: (request, context) => ({
            params: request.params,
            context,
          }),
          getView: () => import('./containers/login'),
        })
  ),
  '/forgot': route({
    title: 'Reset Password',
    getView: () => import('./containers/forgot'),
  }),
  '/recover/:token': route({
    title: 'Recover Password',
    getData: ({ params }) => ({ params }),
    getView: () => import('./containers/reset-password'),
  }),
  '/confirm/:token': route({
    title: 'Confirmation',
    getData: ({ params }) => ({ params }),
    getView: () => import('./containers/confirm'),
  }),
  '/projects': compose(
    withData((request, context) => ({
      params: request.params,
      context,
    })),
    mount({
      '*': map(async (request, context) =>
        !context.currentUser
          ? redirect(
              `/login${
                request.path
                  ? `?redirectTo=${encodeURIComponent(
                      request.mountpath + request.search
                    )}`
                  : ''
              }`
            )
          : mount({
              '/': route({
                title: 'Projects',
                getView: () => import('./containers/projects'),
              }),
              '/create': route({
                title: 'Create Project',
                getView: () => import('./containers/create-project'),
              }),
            })
      ),
    })
  ),

  '/:project': compose(
    withData((request, context) => ({ params: request.params, context })),
    mount({
      '*': map(async (request, context) => {
        if (!context.currentUser) {
          redirect(
            `/login${
              request.path
                ? `?redirectTo=${encodeURIComponent(
                    request.mountpath + request.search
                  )}`
                : ''
            }`
          );
        }
        return mount({
          '/': redirect('devices'),
          '/devices': mount({
            '/': route({
              title: 'Devices',
              getView: () => import('./containers/devices'),
            }),
            '/register': route({
              title: 'Register Device',
              getView: () => import('./containers/register-device'),
            }),
            '/:device': compose(
              withView(() => import('./containers/device/index')),
              withData((request, context) => ({
                params: request.params,
                context,
              })),
              mount({
                '/': redirect('overview'),
                '/overview': route({
                  title: 'Overview - Device',
                  getView: () => import('./containers/device/overview'),
                }),
                '/settings': route({
                  title: 'Settings - Device',
                  getView: () => import('./containers/device/settings'),
                }),
              })
            ),
          }),
          '/iam': compose(
            withView(() => import('./containers/iam')),
            mount({
              '/': redirect('members'),
              '/members': mount({
                '/': route({
                  title: 'Members',
                  getView: () => import('./containers/iam/members'),
                }),
                '/:user': route({
                  title: 'Member',
                  getData: (request, context) => ({
                    params: request.params,
                    context,
                  }),
                  getView: () => import('./containers/iam/member'),
                }),
                '/add': route({
                  title: 'Add Member',
                  getView: () => import('./containers/iam/add-member'),
                }),
              }),
              '/roles': mount({
                '/': route({
                  title: 'Roles',
                  getView: () => import('./containers/iam/roles'),
                }),
                '/:role': route({
                  title: 'Role',
                  getData: (request, context) => ({
                    params: request.params,
                    context,
                  }),
                  getView: () => import('./containers/iam/role'),
                }),
                '/create': route({
                  title: 'Create Role',
                  getView: () => import('./containers/iam/create-role'),
                }),
              }),
              '/service-accounts': mount({
                '/': route({
                  title: 'Service Accounts',
                  getView: () => import('./containers/iam/service-accounts'),
                }),
                '/:service': route({
                  title: 'Service Account',
                  getData: (request, context) => ({
                    params: request.params,
                    context,
                  }),
                  getView: () => import('./containers/iam/service-account'),
                }),
                '/create': route({
                  title: 'Create Service Account',
                  getView: () =>
                    import('./containers/iam/create-service-account'),
                }),
              }),
            })
          ),
          '/applications': mount({
            '/': route({
              title: 'Applications',
              getView: () => import('./containers/applications'),
            }),
            '/create': route({
              title: 'Create Application',
              getView: () => import('./containers/create-application'),
            }),
            '/:application': compose(
              withView(() => import('./containers/application')),
              withData((request, context) => ({
                params: request.params,
                context,
              })),
              mount({
                '/': redirect('overview'),
                '/overview': route({
                  title: 'Overview - Application',
                  getView: () => import('./containers/application/overview'),
                }),
                '/releases': mount({
                  '/': route({
                    title: 'Releases - Application',
                    getView: () => import('./containers/application/releases'),
                  }),
                  '/create': route({
                    title: 'Create Release - Application',
                    getView: () =>
                      import('./containers/application/create-release'),
                  }),
                  '/:release': route({
                    title: 'Release - Application',
                    getData: (request, context) => ({
                      params: request.params,
                      context,
                    }),
                    getView: () => import('./containers/application/release'),
                  }),
                }),
                '/scheduling': route({
                  title: 'Scheduling - Application',
                  getView: () => import('./containers/application/scheduling'),
                }),
                '/release-pinning': route({
                  title: 'Release Pinning - Application',
                  getView: () =>
                    import('./containers/application/release-pinning'),
                }),
                '/settings': route({
                  title: 'Settings - Application',
                  getView: () => import('./containers/application/settings'),
                }),
              })
            ),
          }),
          '/provisioning': mount({
            '/': route({
              title: 'Provisioning',
              getView: () => import('./containers/provisioning'),
            }),
            '/registration-tokens': mount({
              '/create': route({
                title: 'Create Registration Token',
                getView: () => import('./containers/create-registration-token'),
              }),
              '/:token': compose(
                withView(() => import('./containers/registration-token')),
                withData((request, context) => ({
                  params: request.params,
                  context,
                })),
                mount({
                  '/': redirect('overview'),
                  '/overview': route({
                    title: 'Overview - Registration Token',
                    getView: () =>
                      import('./containers/registration-token/overview'),
                  }),
                  '/settings': route({
                    title: 'Settings - Registration Token',
                    getView: () =>
                      import('./containers/registration-token/settings'),
                  }),
                })
              ),
            }),
          }),
          '/monitoring': compose(
            withView(() => import('./containers/monitoring')),
            withData(async ({ params: { project: projectId } }) => {
              const { data: project } = await api.project({
                projectId,
              });
              const { data: applications } = await api.applications({
                projectId,
              });
              const { data: devices } = await api.devices({ projectId });
              return {
                project,
                applications,
                devices,
              };
            }),
            mount({
              '/': route({
                view: ({ route: { data } }) => {
                  const navigation = useNavigation();
                  useEffect(() => {
                    if (data.project.datadogApiKey) {
                      navigation.navigate('monitoring/project');
                    } else {
                      navigation.navigate('monitoring/integrations');
                    }
                  }, []);
                  return null;
                },
              }),
              '/integrations': route({
                title: 'Integrations - Monitoring',
                getView: () => import('./containers/monitoring/integrations'),
              }),
              '/project': route({
                title: 'Project - Monitoring',
                getData: async request => {
                  const {
                    data: { exposedMetrics: metrics },
                  } = await api.projectMetricsConfig({
                    projectId: request.params.project,
                  });

                  return {
                    metrics,
                  };
                },
                getView: () => import('./containers/monitoring/project'),
              }),
              '/device': compose(
                withData(async request => {
                  const {
                    data: { exposedMetrics: metrics },
                  } = await api.deviceMetricsConfig({
                    projectId: request.params.project,
                  });

                  return {
                    metrics,
                  };
                }),
                mount({
                  '/': route({
                    title: 'Device - Monitoring',
                    getView: () => import('./containers/monitoring/device'),
                  }),
                })
              ),
              '/service': compose(
                withData(async request => {
                  const { data: metrics } = await api.serviceMetricsConfig({
                    projectId: request.params.project,
                  });

                  return {
                    metrics,
                  };
                }),
                mount({
                  '/': route({
                    title: 'Service - Monitoring',
                    getView: () => import('./containers/monitoring/service'),
                  }),
                })
              ),
            })
          ),
          '/settings': route({
            title: 'Settings - Project',
            getView: () => import('./containers/project-settings'),
          }),
          '/ssh': route({
            title: 'SSH',
            getData: async request => {
              const response = await api.devices({
                projectId: request.params.project,
              });

              return {
                devices: response.data,
                params: request.params,
              };
            },
            getView: () => import('./containers/ssh'),
          }),
        });
      }),
    })
  ),
});
