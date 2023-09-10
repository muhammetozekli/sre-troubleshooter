package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := signals.SetupSignalHandler()
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		//namespace := "platform"

		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
		if err != nil {
			fmt.Println(err)
		}
		for _, namespace := range namespaces.Items {

			pods, err := clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), v1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			for _, pod := range pods.Items {
				//fmt.Println(pod.Name, pod.Status.Reason)
				if pod.Status.Reason == "Evicted" {
					//fmt.Println(pod.Name+"\n", pod.Status.Reason+"\n", pod.Status.Message)
					err = clientset.CoreV1().Pods(namespace.Name).Delete(context.TODO(), pod.Name, v1.DeleteOptions{})
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(pod.Name + "   \n Deleted")
					}
				} else if pod.Status.Phase != "Running" {
					fmt.Println(namespace.Name, pod.Name+"   \n")
					fmt.Println(len(pod.Status.Conditions))
					if pod.Status.Conditions[0].Type == "Initialized" && pod.Status.Conditions[0].Status == "False" {
						fmt.Println(pod.Name+"  couldnt initialized", pod.Status.Conditions[0].Reason, pod.Status.Conditions[0].Message)
					}
					/*
						if pod.Status.Conditions[1].Type == "Ready" && pod.Status.Conditions[1].Status == "False" {
							fmt.Println(pod.Name, "is not ready", pod.Status.Conditions[1].Reason)
						}
						if pod.Status.Conditions[2].Type == "ContainersNotReady" && pod.Status.Conditions[2].Status == "False" {
							fmt.Println(pod.Name, "is not ready", pod.Status.Conditions[2].Reason, pod.Status.Conditions[2].Message)
						}
						if pod.Status.Conditions[3].Type == "PodScheduled" && pod.Status.Conditions[3].Status == "False" {
							fmt.Println(pod.Name, "couldnt scheduled", pod.Status.Conditions[3].Reason, pod.Status.Conditions[3].Message)
						}

					*/
				}

			}

		}

		//pod := "cloud-cost-exporter-75b5bb6666-l7zjc"
		//_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, v1.GetOptions{})

		/*
			if errors.IsNotFound(err) {
				fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %s in namespace %s: %v\n",
					pod, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
			}

		*/

		y, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})

		for _, node := range y.Items {
			//fmt.Printf("%s\n", node.Name)

			if node.Status.Conditions[4].Type == "Ready" && node.Status.Conditions[4].Status == "False" {
				fmt.Println(node.Name)
			}

			/*
				for _, condition := range node.Status.Conditions {
					fmt.Printf("\t%s: %s\n", condition.Type, condition.Status)

				}
			*/
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
