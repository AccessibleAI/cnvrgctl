cluster_name: cnvrg
ingress:
  provider: none
monitoring:
  provider: none
nodes:
  - address: {{ .Data.Server }}
    user: {{ .Data.User }}
    ssh_key_path: {{ .Data.SshPrivateKey }}
    role:
      - controlplane
      - etcd
      - worker
addons: |-
  ---
  apiVersion: apiextensions.k8s.io/v1
  kind: CustomResourceDefinition
  metadata:
    name: cnvrgapps.mlops.cnvrg.io
  spec:
    group: mlops.cnvrg.io
    names:
      kind: CnvrgApp
      listKind: CnvrgAppList
      plural: cnvrgapps
      singular: cnvrgapp
    scope: Namespaced
    versions:
      - name: v1
        additionalPrinterColumns:
          - description: cnvrg version
            jsonPath: .spec.cnvrgApp.image
            name: Version
            type: string
          - description: otags
            jsonPath: .spec.otags
            name: Otag
            type: string
          - description: cnvrg status
            jsonPath: .status.conditions[0].message
            name: Status
            type: string
          - jsonPath: .metadata.creationTimestamp
            name: Age
            type: date
        schema:
          openAPIV3Schema:
            description: CnvrgApp is the Schema for the cnvrgapps API
            properties:
              apiVersion:
                description: 'APIVersion defines the versioned schema of this representation
                of an object. Servers should convert recognized schemas to the latest
                internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                type: string
              kind:
                description: 'Kind is a string value representing the REST resource this
                object represents. Servers may infer this from the endpoint the client
                submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                type: string
              metadata:
                type: object
              spec:
                description: Spec defines the desired state of CnvrgApp
                type: object
                x-kubernetes-preserve-unknown-fields: true
              status:
                description: Status defines the observed state of CnvrgApp
                type: object
                x-kubernetes-preserve-unknown-fields: true
            type: object
        served: true
        storage: true
        subresources:
          status: {}
  ---
  apiVersion: v1
  kind: Namespace
  metadata:
    name: cnvrg
  ---
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    labels:
      control-plane: cnvrg-operator
    name: cnvrg-operator
    namespace: cnvrg
  spec:
    replicas: 1
    selector:
      matchLabels:
        control-plane: cnvrg-operator
    template:
      metadata:
        labels:
          control-plane: cnvrg-operator
      spec:
        serviceAccountName: cnvrg-operator
        terminationGracePeriodSeconds: 10
        tolerations:
          - key: "cnvrg-taint"
            operator: "Equal"
            value: "true"
            effect: "NoSchedule"
        containers:
          - name: manager
            image: "docker.io/cnvrg/cnvrg-operator:2.21.0"
            args:
              - --enable-leader-election
              - --leader-election-id=cnvrg-operator
              - --leader-election-namespace=cnvrg
              - --max-concurrent-reconciles=8
              - --zap-encoder=console
              - --zap-log-level=info
            env:
              - name: ANSIBLE_JINJA2_NATIVE
                value: "true"
              - name: ANSIBLE_HASH_BEHAVIOUR
                value: merge
              - name: WATCH_NAMESPACE
                valueFrom:
                  fieldRef:
                    fieldPath: metadata.namespace
              - name: POD_NAME
                valueFrom:
                  fieldRef:
                    fieldPath: metadata.name
  ---
  apiVersion: rbac.authorization.k8s.io/v1
  kind: Role
  metadata:
    name: cnvrg-operator
    namespace: cnvrg
  rules:
    - apiGroups:
        - batch
      resources:
        - cronjobs
      verbs:
        - '*'
    - apiGroups:
        - networking.k8s.io
      resources:
        - ingresses
      verbs:
        - '*'
    - apiGroups:
        - route.openshift.io
      resources:
        - '*'
      verbs:
        - '*'
    - apiGroups:
        - rbac.authorization.k8s.io
      resources:
        - roles
        - rolebindings
      verbs:
        - "*"
    - apiGroups:
        - ""
      resources:
        - pods
        - services
        - services/finalizers
        - endpoints
        - persistentvolumeclaims
        - events
        - configmaps
        - secrets
        - serviceaccounts
      verbs:
        - create
        - delete
        - get
        - list
        - patch
        - update
        - watch
    - apiGroups:
        - batch
        - extensions
        - cronjobs
      resources:
        - jobs
      verbs:
        - create
        - delete
        - get
        - list
        - patch
        - update
        - watch
    - apiGroups:
        - apps
      resources:
        - deployments
        - daemonsets
        - replicasets
        - statefulsets
      verbs:
        - create
        - delete
        - get
        - list
        - patch
        - update
        - watch
    - apiGroups:
        - monitoring.coreos.com
      resources:
        - servicemonitors
        - prometheuses
        - prometheusrules
      verbs:
        - list
        - get
        - create
        - patch
        - watch
    - apiGroups:
        - apps
      resourceNames:
        - cnvrg-operator
      resources:
        - deployments/finalizers
      verbs:
        - update
    - apiGroups:
        - ""
      resources:
        - pods
      verbs:
        - get
    - apiGroups:
        - apps
      resources:
        - replicasets
        - deployments
      verbs:
        - get
    - apiGroups:
        - mlops.cnvrg.io
      resources:
        - '*'
      verbs:
        - '*'
  ---
  kind: RoleBinding
  apiVersion: rbac.authorization.k8s.io/v1
  metadata:
    name: cnvrg-operator
    namespace: cnvrg
  subjects:
    - kind: ServiceAccount
      name: cnvrg-operator
  roleRef:
    kind: Role
    name: cnvrg-operator
    apiGroup: rbac.authorization.k8s.io
  ---
  apiVersion: v1
  kind: ServiceAccount
  metadata:
    name: cnvrg-operator
    namespace: cnvrg